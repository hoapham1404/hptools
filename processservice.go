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

type ProcessService struct{}

// ProcessInfo holds a single row from `tasklist`
type ProcessInfo struct {
	ImageName   string
	PID         int
	SessionName string
	SessionNum  int
	MemUsageB   uint64 // bytes
	MemUsageStr string // original string, e.g. "12,345 K"
}

// ListProcessesWindow returns structured data parsed from `tasklist` (Windows only)
func (p *ProcessService) ListProcessesWindow() ([]ProcessInfo, error) {
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

	procs := make([]ProcessInfo, 0, len(rows))
	for _, rec := range rows {
		// Expect: Image Name, PID, Session Name, Session#, Mem Usage
		if len(rec) < 5 {
			continue
		}
		pid, _ := strconv.Atoi(strings.ReplaceAll(rec[1], ",", ""))
		sessNum, _ := strconv.Atoi(strings.ReplaceAll(rec[3], ",", ""))
		memStr := rec[4]
		procs = append(procs, ProcessInfo{
			ImageName:   rec[0],
			PID:         pid,
			SessionName: rec[2],
			SessionNum:  sessNum,
			MemUsageB:   parseMemBytes(memStr),
			MemUsageStr: memStr,
		})
	}
	return procs, nil
}

func parseMemBytes(s string) uint64 {
	// e.g. "12,345 K" -> 12641880 bytes (12345 * 1024)
	s = strings.TrimSpace(s)
	var digits strings.Builder
	for _, ch := range s {
		if ch >= '0' && ch <= '9' {
			digits.WriteRune(ch)
		}
	}
	n, _ := strconv.ParseUint(digits.String(), 10, 64)
	return n * 1024
}

// Add Windows API declarations
var (
	user32                       = syscall.NewLazyDLL("user32.dll")
	procFindWindow               = user32.NewProc("FindWindowW")
	procSetWindowPos             = user32.NewProc("SetWindowPos")
	procGetWindowRect            = user32.NewProc("GetWindowRect")
	procEnumWindows              = user32.NewProc("EnumWindows")
	procGetWindowThreadProcessId = user32.NewProc("GetWindowThreadProcessId")
)

// RECT structure for window coordinates
type RECT struct {
	Left, Top, Right, Bottom int32
}

// SetWindowSize sets the size of a window by process PID
func (p *ProcessService) SetWindowSize(pid int, width, height int) error {
	hwnd, err := p.findWindowByPID(pid)
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

// SetWindowPosition sets the position and size of a window
func (p *ProcessService) SetWindowPosition(pid int, x, y, width, height int) error {
	hwnd, err := p.findWindowByPID(pid)
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

// GetWindowSize gets the current size and position of a window
func (p *ProcessService) GetWindowSize(pid int) (x, y, width, height int, err error) {
	hwnd, err := p.findWindowByPID(pid)
	if err != nil {
		return 0, 0, 0, 0, fmt.Errorf("finding window for PID %d: %w", pid, err)
	}

	var rect RECT
	ret, _, _ := procGetWindowRect.Call(uintptr(hwnd), uintptr(unsafe.Pointer(&rect)))

	if ret == 0 {
		return 0, 0, 0, 0, fmt.Errorf("GetWindowRect failed")
	}

	return int(rect.Left), int(rect.Top),
		int(rect.Right - rect.Left), int(rect.Bottom - rect.Top), nil
}

// Helper function to find window handle by PID
func (p *ProcessService) findWindowByPID(targetPID int) (syscall.Handle, error) {
	var foundHwnd syscall.Handle

	// Callback function for EnumWindows
	enumProc := syscall.NewCallback(func(hwnd syscall.Handle, lParam uintptr) uintptr {
		var pid uint32
		procGetWindowThreadProcessId.Call(uintptr(hwnd), uintptr(unsafe.Pointer(&pid)))

		if int(pid) == targetPID {
			foundHwnd = hwnd
			return 0 // Stop enumeration
		}
		return 1 // Continue enumeration
	})

	procEnumWindows.Call(enumProc, 0)

	if foundHwnd == 0 {
		return 0, fmt.Errorf("no window found for PID %d", targetPID)
	}

	return foundHwnd, nil
}

// GetApplicationProcesses returns only processes that are likely applications
func (p *ProcessService) GetApplicationProcesses() ([]ProcessInfo, error) {
	allProcs, err := p.ListProcessesWindow()
	if err != nil {
		return nil, err
	}

	var apps []ProcessInfo
	for _, proc := range allProcs {
		if p.isApplication(proc) {
			apps = append(apps, proc)
		}
	}

	return apps, nil
}

// isApplication determines if a process is likely an application
func (p *ProcessService) isApplication(proc ProcessInfo) bool {
	imageName := strings.ToLower(proc.ImageName)

	// System processes to exclude
	systemProcesses := []string{
		"system", "smss.exe", "csrss.exe", "wininit.exe", "winlogon.exe",
		"services.exe", "lsass.exe", "svchost.exe", "spoolsv.exe",
		"dwm.exe", "audiodg.exe", "conhost.exe", "taskmgr.exe",
		"cmd.exe", "powershell.exe", "wuauclt.exe",
	}

	// Check if it's a system process
	for _, sys := range systemProcesses {
		if imageName == sys {
			return false
		}
	}

	// Applications typically run in user sessions (not session 0)
	// and have .exe extension
	return proc.SessionNum > 0 && strings.HasSuffix(imageName, ".exe")
}
