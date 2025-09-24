package services

import (
	"hptools/internal/models"
)

// WailsWindowService is the concrete implementation for Wails
type WailsWindowService struct {
	service WindowService
}

// NewWailsWindowService creates a new Wails-compatible service
func NewWailsWindowService(service WindowService) *WailsWindowService {
	return &WailsWindowService{service: service}
}

// GetApplicationProcesses returns only processes that have visible windows
func (w *WailsWindowService) GetApplicationProcesses() ([]models.ProcessInfo, error) {
	return w.service.GetApplicationProcesses()
}

// GetAllProcessesWithWindows returns all processes that have visible windows (for debugging)
func (w *WailsWindowService) GetAllProcessesWithWindows() ([]models.ProcessInfo, error) {
	return w.service.GetAllProcessesWithWindows()
}

// SetWindowSize sets the size of a window by process PID, keeping current position
func (w *WailsWindowService) SetWindowSize(pid int, width, height int) error {
	return w.service.SetWindowSize(pid, width, height)
}

// SetWindowPosition sets both position and size of a window
func (w *WailsWindowService) SetWindowPosition(pid int, x, y, width, height int) error {
	return w.service.SetWindowPosition(pid, x, y, width, height)
}

// GetWindowInfo gets the current size and position of a window
func (w *WailsWindowService) GetWindowInfo(pid int) (*models.WindowInfo, error) {
	return w.service.GetWindowInfo(pid)
}
