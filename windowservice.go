package main

import (
	"encoding/csv"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
	"unsafe"
)

// WindowService provides window management functionality
type WindowService struct{}

// ProcessInfo represents information about a running process
type ProcessInfo struct {
	ImageName   string `json:"imageName"`
	PID         int    `json:"pid"`
	SessionName string `json:"sessionName"`
	SessionNum  int    `json:"sessionNum"`
	MemUsageB   int64  `json:"memUsageB"`
	MemUsageStr string `json:"memUsageStr"`
	WindowTitle string `json:"windowTitle"`
	HasWindow   bool   `json:"hasWindow"`
	WindowCount int    `json:"windowCount"`
}

// WindowInfo represents window position and size information
type WindowInfo struct {
	X      int `json:"x"`
	Y      int `json:"y"`
	Width  int `json:"width"`
	Height int `json:"height"`
}

// Windows API declarations
var (
	user32                       = syscall.NewLazyDLL("user32.dll")
	procFindWindow               = user32.NewProc("FindWindowW")
	procSetWindowPos             = user32.NewProc("SetWindowPos")
	procGetWindowRect            = user32.NewProc("GetWindowRect")
	procEnumWindows              = user32.NewProc("EnumWindows")
	procGetWindowThreadProcessId = user32.NewProc("GetWindowThreadProcessId")
	procIsWindowVisible          = user32.NewProc("IsWindowVisible")
	procGetWindowTextW           = user32.NewProc("GetWindowTextW")
)

// RECT structure for window coordinates
type RECT struct {
	Left, Top, Right, Bottom int32
}

// GetApplicationProcesses returns only processes that have visible windows
func (w *WindowService) GetApplicationProcesses() ([]ProcessInfo, error) {
	out, err := exec.Command("tasklist", "/FO", "CSV", "/NH").Output()
	if err != nil {
		return nil, fmt.Errorf("executing tasklist: %w", err)
	}

	r := csv.NewReader(strings.NewReader(string(out)))
	r.FieldsPerRecord = -1

	rows, err := r.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("parsing CSV: %w", err)
	}

	var apps []ProcessInfo
	for _, rec := range rows {
		if len(rec) < 5 {
			continue
		}

		proc := ProcessInfo{
			ImageName:   strings.Trim(rec[0], "\""),
			SessionName: strings.Trim(rec[2], "\""),
			MemUsageStr: strings.Trim(rec[4], "\""),
		}

		// Parse PID
		pidStr := strings.ReplaceAll(strings.Trim(rec[1], "\""), ",", "")
		proc.PID, _ = strconv.Atoi(pidStr)

		// Parse Session Number
		sessStr := strings.ReplaceAll(strings.Trim(rec[3], "\""), ",", "")
		proc.SessionNum, _ = strconv.Atoi(sessStr)

		// Parse memory usage
		proc.MemUsageB = parseMemBytes(proc.MemUsageStr)

		// Filter to only include applications (not system processes)
		if w.isApplication(proc) {
			// Check if this process actually has visible windows
			windowInfo := w.getProcessWindowInfo(proc.PID)
			if windowInfo.HasWindow {
				proc.WindowTitle = windowInfo.WindowTitle
				proc.HasWindow = windowInfo.HasWindow
				proc.WindowCount = windowInfo.WindowCount
				apps = append(apps, proc)
			}
		}
	}

	// Group processes by image name and only keep the one with the main window
	groupedApps := w.groupAndFilterProcesses(apps)

	return groupedApps, nil
}

// ProcessWindowInfo contains window information for a process
type ProcessWindowInfo struct {
	HasWindow     bool
	WindowTitle   string
	WindowCount   int
	MainWindowPID int
}

// getProcessWindowInfo checks if a process has visible windows and gets window info
func (w *WindowService) getProcessWindowInfo(targetPID int) ProcessWindowInfo {
	var windowInfo ProcessWindowInfo
	var foundWindows []string

	// Callback function for EnumWindows
	enumProc := syscall.NewCallback(func(hwnd syscall.Handle, lParam uintptr) uintptr {
		// Check if window is visible
		if ret, _, _ := procIsWindowVisible.Call(uintptr(hwnd)); ret == 0 {
			return 1 // Continue enumeration for invisible windows
		}

		var pid uint32
		procGetWindowThreadProcessId.Call(uintptr(hwnd), uintptr(unsafe.Pointer(&pid)))

		if int(pid) == targetPID {
			// Get window title
			title := w.getWindowTitle(hwnd)
			if title != "" && title != "Default IME" && title != "MSCTFIME UI" {
				foundWindows = append(foundWindows, title)
				if windowInfo.WindowTitle == "" || len(title) > len(windowInfo.WindowTitle) {
					windowInfo.WindowTitle = title // Keep the longest/most descriptive title
				}
			}
			windowInfo.WindowCount++
		}
		return 1 // Continue enumeration
	})

	procEnumWindows.Call(enumProc, 0)

	windowInfo.HasWindow = windowInfo.WindowCount > 0 && windowInfo.WindowTitle != ""
	return windowInfo
}

// getWindowTitle retrieves the title of a window
func (w *WindowService) getWindowTitle(hwnd syscall.Handle) string {
	const maxLength = 256
	buf := make([]uint16, maxLength)
	ret, _, _ := procGetWindowTextW.Call(uintptr(hwnd), uintptr(unsafe.Pointer(&buf[0])), maxLength)
	if ret > 0 {
		return syscall.UTF16ToString(buf[:ret])
	}
	return ""
}

// groupAndFilterProcesses groups processes by image name and keeps only the main process with window
func (w *WindowService) groupAndFilterProcesses(processes []ProcessInfo) []ProcessInfo {
	processGroups := make(map[string][]ProcessInfo)

	// Group by image name
	for _, proc := range processes {
		key := strings.ToLower(proc.ImageName)
		processGroups[key] = append(processGroups[key], proc)
	}

	var result []ProcessInfo
	for _, group := range processGroups {
		if len(group) == 1 {
			// Single process, keep it
			result = append(result, group[0])
		} else {
			// Multiple processes, find the best candidate
			best := w.selectBestProcess(group)
			if best != nil {
				result = append(result, *best)
			}
		}
	}

	return result
}

// selectBestProcess selects the best process from a group of processes with the same name
func (w *WindowService) selectBestProcess(processes []ProcessInfo) *ProcessInfo {
	var best *ProcessInfo

	for i := range processes {
		proc := &processes[i]

		if best == nil {
			best = proc
			continue
		}

		// Prefer process with meaningful window title
		if len(proc.WindowTitle) > len(best.WindowTitle) {
			best = proc
			continue
		}

		// Prefer process with higher memory usage (likely the main process)
		if proc.MemUsageB > best.MemUsageB {
			best = proc
			continue
		}

		// Prefer lower PID (started earlier, likely parent process)
		if proc.PID < best.PID {
			best = proc
		}
	}

	return best
}

// SetWindowSize sets the size of a window by process PID, keeping current position
func (w *WindowService) SetWindowSize(pid int, width, height int) error {
	hwnd, err := w.findWindowByPID(pid)
	if err != nil {
		return fmt.Errorf("finding window for PID %d: %w", pid, err)
	}

	// SWP_NOMOVE (0x0002) - retain current position, only change size
	// SWP_NOZORDER (0x0004) - retain current Z order
	// SWP_NOACTIVATE (0x0010) - do not activate window
	const SWP_NOMOVE = 0x0002
	const SWP_NOZORDER = 0x0004
	const SWP_NOACTIVATE = 0x0010

	ret, _, _ := procSetWindowPos.Call(
		uintptr(hwnd),
		0, // hWndInsertAfter (ignored due to SWP_NOZORDER)
		0, // X (ignored due to SWP_NOMOVE)
		0, // Y (ignored due to SWP_NOMOVE)
		uintptr(width),
		uintptr(height),
		SWP_NOMOVE|SWP_NOZORDER|SWP_NOACTIVATE,
	)

	if ret == 0 {
		return fmt.Errorf("SetWindowPos failed")
	}

	return nil
}

// SetWindowPosition sets both position and size of a window
func (w *WindowService) SetWindowPosition(pid int, x, y, width, height int) error {
	hwnd, err := w.findWindowByPID(pid)
	if err != nil {
		return fmt.Errorf("finding window for PID %d: %w", pid, err)
	}

	const SWP_NOZORDER = 0x0004
	const SWP_NOACTIVATE = 0x0010

	ret, _, _ := procSetWindowPos.Call(
		uintptr(hwnd),
		0, // hWndInsertAfter
		uintptr(x),
		uintptr(y),
		uintptr(width),
		uintptr(height),
		SWP_NOZORDER|SWP_NOACTIVATE,
	)

	if ret == 0 {
		return fmt.Errorf("SetWindowPos failed")
	}

	return nil
}

// GetWindowInfo gets the current size and position of a window
func (w *WindowService) GetWindowInfo(pid int) (*WindowInfo, error) {
	hwnd, err := w.findWindowByPID(pid)
	if err != nil {
		return nil, fmt.Errorf("finding window for PID %d: %w", pid, err)
	}

	var rect RECT
	ret, _, _ := procGetWindowRect.Call(uintptr(hwnd), uintptr(unsafe.Pointer(&rect)))

	if ret == 0 {
		return nil, fmt.Errorf("GetWindowRect failed")
	}

	return &WindowInfo{
		X:      int(rect.Left),
		Y:      int(rect.Top),
		Width:  int(rect.Right - rect.Left),
		Height: int(rect.Bottom - rect.Top),
	}, nil
}

// Helper function to find window handle by PID
func (w *WindowService) findWindowByPID(targetPID int) (syscall.Handle, error) {
	var foundWindows []syscall.Handle

	// Callback function for EnumWindows
	enumProc := syscall.NewCallback(func(hwnd syscall.Handle, lParam uintptr) uintptr {
		// Check if window is visible
		if ret, _, _ := procIsWindowVisible.Call(uintptr(hwnd)); ret == 0 {
			return 1 // Continue enumeration for invisible windows
		}

		var pid uint32
		procGetWindowThreadProcessId.Call(uintptr(hwnd), uintptr(unsafe.Pointer(&pid)))

		if int(pid) == targetPID {
			title := w.getWindowTitle(hwnd)
			// Filter out system windows and empty titles
			if title != "" && title != "Default IME" && title != "MSCTFIME UI" &&
				!strings.Contains(title, "Program Manager") {
				foundWindows = append(foundWindows, hwnd)
			}
		}
		return 1 // Continue enumeration
	})

	procEnumWindows.Call(enumProc, 0)

	if len(foundWindows) == 0 {
		return 0, fmt.Errorf("no visible window found for PID %d", targetPID)
	}

	// If multiple windows, try to find the main window
	if len(foundWindows) > 1 {
		for _, hwnd := range foundWindows {
			title := w.getWindowTitle(hwnd)
			// Prefer windows with meaningful titles (not just filenames)
			if len(title) > 10 && !strings.HasSuffix(strings.ToLower(title), ".exe") {
				return hwnd, nil
			}
		}
	}

	// Return the first window found
	return foundWindows[0], nil
}

// GetAllProcessesWithWindows returns all processes that have visible windows (for debugging)
func (w *WindowService) GetAllProcessesWithWindows() ([]ProcessInfo, error) {
	out, err := exec.Command("tasklist", "/FO", "CSV", "/NH").Output()
	if err != nil {
		return nil, fmt.Errorf("executing tasklist: %w", err)
	}

	r := csv.NewReader(strings.NewReader(string(out)))
	r.FieldsPerRecord = -1

	rows, err := r.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("parsing CSV: %w", err)
	}

	var apps []ProcessInfo
	for _, rec := range rows {
		if len(rec) < 5 {
			continue
		}

		proc := ProcessInfo{
			ImageName:   strings.Trim(rec[0], "\""),
			SessionName: strings.Trim(rec[2], "\""),
			MemUsageStr: strings.Trim(rec[4], "\""),
		}

		// Parse PID
		pidStr := strings.ReplaceAll(strings.Trim(rec[1], "\""), ",", "")
		proc.PID, _ = strconv.Atoi(pidStr)

		// Parse Session Number
		sessStr := strings.ReplaceAll(strings.Trim(rec[3], "\""), ",", "")
		proc.SessionNum, _ = strconv.Atoi(sessStr)

		// Parse memory usage
		proc.MemUsageB = parseMemBytes(proc.MemUsageStr)

		// Check if this process has visible windows (no application filtering)
		windowInfo := w.getProcessWindowInfo(proc.PID)
		if windowInfo.HasWindow {
			proc.WindowTitle = windowInfo.WindowTitle
			proc.HasWindow = windowInfo.HasWindow
			proc.WindowCount = windowInfo.WindowCount
			apps = append(apps, proc)
		}
	}

	return apps, nil
}

// isApplication determines if a process is likely an application (not a system process)
func (w *WindowService) isApplication(proc ProcessInfo) bool {
	imageName := strings.ToLower(proc.ImageName)

	// System processes to exclude
	systemProcesses := []string{
		"system", "smss.exe", "csrss.exe", "wininit.exe", "winlogon.exe",
		"services.exe", "lsass.exe", "svchost.exe", "spoolsv.exe",
		"dwm.exe", "audiodg.exe", "conhost.exe", "taskmgr.exe",
		"cmd.exe", "powershell.exe", "wuauclt.exe", "mmc.exe",
		"rundll32.exe", "dllhost.exe", "sihost.exe", "fontdrvhost.exe",
		"winrt.exe", "backgroundtaskhost.exe", "runtimebroker.exe",
	}

	// Check if it's a known system process
	for _, sys := range systemProcesses {
		if imageName == sys {
			return false
		}
	}

	// Applications typically run in user sessions (not session 0)
	// and have .exe extension
	return proc.SessionNum > 0 && strings.HasSuffix(imageName, ".exe")
}

// parseMemBytes converts memory usage string to bytes
func parseMemBytes(memStr string) int64 {
	// Remove commas and "K" suffix, convert to bytes
	memStr = strings.ReplaceAll(memStr, ",", "")
	memStr = strings.TrimSuffix(strings.TrimSpace(memStr), " K")

	if val, err := strconv.ParseInt(memStr, 10, 64); err == nil {
		return val * 1024 // Convert KB to bytes
	}
	return 0
}
