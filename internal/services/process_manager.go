package services

import (
	"encoding/csv"
	"fmt"
	"log/slog"
	"os/exec"
	"strconv"
	"strings"
	"syscall"

	"hptools/internal/models"
	"hptools/internal/windows"
)

type processManager struct {
	api    *windows.API
	logger *slog.Logger
}

// NewProcessManager creates a new process manager
func NewProcessManager(api *windows.API, logger *slog.Logger) ProcessManager {
	return &processManager{
		api:    api,
		logger: logger,
	}
}

// GetApplicationProcesses returns only processes that have visible windows
func (p *processManager) GetApplicationProcesses() ([]models.ProcessInfo, error) {
	processes, err := p.getAllProcesses()
	if err != nil {
		return nil, fmt.Errorf("getting all processes: %w", err)
	}

	var apps []models.ProcessInfo
	for _, proc := range processes {
		if p.IsApplication(proc) {
			// Check if this process actually has visible windows
			windowInfo := p.getProcessWindowInfo(proc.PID)
			if windowInfo.HasWindow {
				proc.WindowTitle = windowInfo.WindowTitle
				proc.HasWindow = windowInfo.HasWindow
				proc.WindowCount = windowInfo.WindowCount
				apps = append(apps, proc)
			}
		}
	}

	// Group processes by image name and only keep the one with the main window
	groupedApps := p.groupAndFilterProcesses(apps)
	p.logger.Info("Found application processes", "count", len(groupedApps))

	return groupedApps, nil
}

// GetAllProcessesWithWindows returns all processes that have visible windows (for debugging)
func (p *processManager) GetAllProcessesWithWindows() ([]models.ProcessInfo, error) {
	processes, err := p.getAllProcesses()
	if err != nil {
		return nil, fmt.Errorf("getting all processes: %w", err)
	}

	var apps []models.ProcessInfo
	for _, proc := range processes {
		// Check if this process has visible windows (no application filtering)
		windowInfo := p.getProcessWindowInfo(proc.PID)
		if windowInfo.HasWindow {
			proc.WindowTitle = windowInfo.WindowTitle
			proc.HasWindow = windowInfo.HasWindow
			proc.WindowCount = windowInfo.WindowCount
			apps = append(apps, proc)
		}
	}

	p.logger.Debug("Found processes with windows", "count", len(apps))
	return apps, nil
}

// IsApplication determines if a process is likely an application (not a system process)
func (p *processManager) IsApplication(proc models.ProcessInfo) bool {
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

// getAllProcesses gets all running processes using tasklist
func (p *processManager) getAllProcesses() ([]models.ProcessInfo, error) {
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

	var processes []models.ProcessInfo
	for _, rec := range rows {
		if len(rec) < 5 {
			continue
		}

		proc := models.ProcessInfo{
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

		processes = append(processes, proc)
	}

	return processes, nil
}

// getProcessWindowInfo checks if a process has visible windows and gets window info
func (p *processManager) getProcessWindowInfo(targetPID int) models.ProcessWindowInfo {
	var windowInfo models.ProcessWindowInfo
	var foundWindows []string

	// Callback function for EnumWindows
	enumProc := syscall.NewCallback(func(hwnd syscall.Handle, lParam uintptr) uintptr {
		// Check if window is visible
		if !p.api.IsWindowVisible(hwnd) {
			return 1 // Continue enumeration for invisible windows
		}

		pid := p.api.GetWindowThreadProcessId(hwnd)
		if int(pid) == targetPID {
			// Get window title
			title := p.api.GetWindowText(hwnd)
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

	if err := p.api.EnumWindows(enumProc); err != nil {
		p.logger.Warn("Failed to enumerate windows", "pid", targetPID, "error", err)
	}

	windowInfo.HasWindow = windowInfo.WindowCount > 0 && windowInfo.WindowTitle != ""
	return windowInfo
}

// groupAndFilterProcesses groups processes by image name and keeps only the main process with window
func (p *processManager) groupAndFilterProcesses(processes []models.ProcessInfo) []models.ProcessInfo {
	processGroups := make(map[string][]models.ProcessInfo)

	// Group by image name
	for _, proc := range processes {
		key := strings.ToLower(proc.ImageName)
		processGroups[key] = append(processGroups[key], proc)
	}

	var result []models.ProcessInfo
	for _, group := range processGroups {
		if len(group) == 1 {
			// Single process, keep it
			result = append(result, group[0])
		} else {
			// Multiple processes, find the best candidate
			best := p.selectBestProcess(group)
			if best != nil {
				result = append(result, *best)
			}
		}
	}

	return result
}

// selectBestProcess selects the best process from a group of processes with the same name
func (p *processManager) selectBestProcess(processes []models.ProcessInfo) *models.ProcessInfo {
	var best *models.ProcessInfo

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
