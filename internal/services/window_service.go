package services

import (
	"log/slog"

	"hptools/internal/windows"
)

// combinedWindowService implements both ProcessManager and WindowManager
type combinedWindowService struct {
	ProcessManager
	WindowManager
}

// NewWindowService creates a new combined window service
func NewWindowService(api *windows.API, logger *slog.Logger) WindowService {
	return &combinedWindowService{
		ProcessManager: NewProcessManager(api, logger),
		WindowManager:  NewWindowManager(api, logger),
	}
}
