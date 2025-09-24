package models

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

// ProcessWindowInfo contains window information for a process
type ProcessWindowInfo struct {
	HasWindow     bool
	WindowTitle   string
	WindowCount   int
	MainWindowPID int
}
