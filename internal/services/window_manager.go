package services

import (
	"fmt"
	"log/slog"
	"strings"
	"syscall"

	"hptools/internal/models"
	"hptools/internal/windows"
)

const errFindingWindowForPID = "finding window for PID %d: %w"

type windowManager struct {
	api    *windows.API
	logger *slog.Logger
}

// NewWindowManager creates a new window manager
func NewWindowManager(api *windows.API, logger *slog.Logger) WindowManager {
	return &windowManager{
		api:    api,
		logger: logger,
	}
}

// SetWindowSize sets the size of a window by process PID, keeping current position
func (w *windowManager) SetWindowSize(pid int, width, height int) error {
	hwnd, err := w.FindWindowByPID(pid)
	if err != nil {
		return fmt.Errorf(errFindingWindowForPID, pid, err)
	}

	err = w.api.SetWindowPos(
		syscall.Handle(hwnd),
		0, 0, width, height,
		windows.SWP_NOMOVE|windows.SWP_NOZORDER|windows.SWP_NOACTIVATE,
	)
	if err != nil {
		return fmt.Errorf("setting window size: %w", err)
	}

	w.logger.Info("Window size changed", "pid", pid, "width", width, "height", height)
	return nil
}

// SetWindowPosition sets both position and size of a window
func (w *windowManager) SetWindowPosition(pid int, x, y, width, height int) error {
	hwnd, err := w.FindWindowByPID(pid)
	if err != nil {
		return fmt.Errorf(errFindingWindowForPID, pid, err)
	}

	err = w.api.SetWindowPos(
		syscall.Handle(hwnd),
		x, y, width, height,
		windows.SWP_NOZORDER|windows.SWP_NOACTIVATE,
	)
	if err != nil {
		return fmt.Errorf("setting window position: %w", err)
	}

	w.logger.Info("Window position changed", "pid", pid, "x", x, "y", y, "width", width, "height", height)
	return nil
}

// GetWindowInfo gets the current size and position of a window
func (w *windowManager) GetWindowInfo(pid int) (*models.WindowInfo, error) {
	hwnd, err := w.FindWindowByPID(pid)
	if err != nil {
		return nil, fmt.Errorf(errFindingWindowForPID, pid, err)
	}

	rect, err := w.api.GetWindowRect(syscall.Handle(hwnd))
	if err != nil {
		return nil, fmt.Errorf("getting window rect: %w", err)
	}

	return &models.WindowInfo{
		X:      int(rect.Left),
		Y:      int(rect.Top),
		Width:  int(rect.Right - rect.Left),
		Height: int(rect.Bottom - rect.Top),
	}, nil
}

// FindWindowByPID finds window handle by PID
func (w *windowManager) FindWindowByPID(targetPID int) (uintptr, error) {
	var foundWindows []syscall.Handle

	// Callback function for EnumWindows
	enumProc := syscall.NewCallback(func(hwnd syscall.Handle, lParam uintptr) uintptr {
		// Check if window is visible
		if !w.api.IsWindowVisible(hwnd) {
			return 1 // Continue enumeration for invisible windows
		}

		pid := w.api.GetWindowThreadProcessId(hwnd)
		if int(pid) == targetPID {
			title := w.api.GetWindowText(hwnd)
			// Filter out system windows and empty titles
			if title != "" && title != "Default IME" && title != "MSCTFIME UI" &&
				!strings.Contains(title, "Program Manager") {
				foundWindows = append(foundWindows, hwnd)
			}
		}
		return 1 // Continue enumeration
	})

	if err := w.api.EnumWindows(enumProc); err != nil {
		return 0, fmt.Errorf("enumerating windows: %w", err)
	}

	if len(foundWindows) == 0 {
		return 0, fmt.Errorf("no visible window found for PID %d", targetPID)
	}

	// If multiple windows, try to find the main window
	if len(foundWindows) > 1 {
		for _, hwnd := range foundWindows {
			title := w.api.GetWindowText(hwnd)
			// Prefer windows with meaningful titles (not just filenames)
			if len(title) > 10 && !strings.HasSuffix(strings.ToLower(title), ".exe") {
				return uintptr(hwnd), nil
			}
		}
	}

	// Return the first window found
	return uintptr(foundWindows[0]), nil
}
