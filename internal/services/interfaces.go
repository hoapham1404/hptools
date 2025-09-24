package services

import "hptools/internal/models"

// ProcessManager defines the interface for process management operations
type ProcessManager interface {
	GetApplicationProcesses() ([]models.ProcessInfo, error)
	GetAllProcessesWithWindows() ([]models.ProcessInfo, error)
	IsApplication(proc models.ProcessInfo) bool
}

// WindowManager defines the interface for window management operations
type WindowManager interface {
	SetWindowSize(pid int, width, height int) error
	SetWindowPosition(pid int, x, y, width, height int) error
	GetWindowInfo(pid int) (*models.WindowInfo, error)
	FindWindowByPID(pid int) (uintptr, error)
}

// WindowService combines both process and window management
type WindowService interface {
	ProcessManager
	WindowManager
}
