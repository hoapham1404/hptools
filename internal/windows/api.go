package windows

import (
	"syscall"
	"unsafe"

	"hptools/internal/models"
)

// API wraps the Windows API calls for better testability and maintainability
type API struct {
	user32                       *syscall.LazyDLL
	procFindWindow               *syscall.LazyProc
	procSetWindowPos             *syscall.LazyProc
	procGetWindowRect            *syscall.LazyProc
	procEnumWindows              *syscall.LazyProc
	procGetWindowThreadProcessId *syscall.LazyProc
	procIsWindowVisible          *syscall.LazyProc
	procGetWindowTextW           *syscall.LazyProc
}

// NewAPI creates a new Windows API wrapper
func NewAPI() *API {
	user32 := syscall.NewLazyDLL("user32.dll")
	return &API{
		user32:                       user32,
		procFindWindow:               user32.NewProc("FindWindowW"),
		procSetWindowPos:             user32.NewProc("SetWindowPos"),
		procGetWindowRect:            user32.NewProc("GetWindowRect"),
		procEnumWindows:              user32.NewProc("EnumWindows"),
		procGetWindowThreadProcessId: user32.NewProc("GetWindowThreadProcessId"),
		procIsWindowVisible:          user32.NewProc("IsWindowVisible"),
		procGetWindowTextW:           user32.NewProc("GetWindowTextW"),
	}
}

// EnumWindows enumerates all windows with the provided callback
func (api *API) EnumWindows(callback uintptr) error {
	ret, _, _ := api.procEnumWindows.Call(callback, 0)
	if ret == 0 {
		return syscall.GetLastError()
	}
	return nil
}

// IsWindowVisible checks if a window is visible
func (api *API) IsWindowVisible(hwnd syscall.Handle) bool {
	ret, _, _ := api.procIsWindowVisible.Call(uintptr(hwnd))
	return ret != 0
}

// GetWindowThreadProcessId gets the process ID for a window
func (api *API) GetWindowThreadProcessId(hwnd syscall.Handle) uint32 {
	var pid uint32
	api.procGetWindowThreadProcessId.Call(uintptr(hwnd), uintptr(unsafe.Pointer(&pid)))
	return pid
}

// GetWindowText gets the window title
func (api *API) GetWindowText(hwnd syscall.Handle) string {
	const maxLength = 256
	buf := make([]uint16, maxLength)
	ret, _, _ := api.procGetWindowTextW.Call(uintptr(hwnd), uintptr(unsafe.Pointer(&buf[0])), maxLength)
	if ret > 0 {
		return syscall.UTF16ToString(buf[:ret])
	}
	return ""
}

// SetWindowPos sets the window position and size
func (api *API) SetWindowPos(hwnd syscall.Handle, x, y, width, height int, flags uint32) error {
	ret, _, _ := api.procSetWindowPos.Call(
		uintptr(hwnd),
		0, // hWndInsertAfter
		uintptr(x),
		uintptr(y),
		uintptr(width),
		uintptr(height),
		uintptr(flags),
	)
	if ret == 0 {
		return syscall.GetLastError()
	}
	return nil
}

// GetWindowRect gets the window rectangle
func (api *API) GetWindowRect(hwnd syscall.Handle) (*models.RECT, error) {
	var rect models.RECT
	ret, _, _ := api.procGetWindowRect.Call(uintptr(hwnd), uintptr(unsafe.Pointer(&rect)))
	if ret == 0 {
		return nil, syscall.GetLastError()
	}
	return &rect, nil
}

// Window position flags
const (
	SWP_NOMOVE     = 0x0002
	SWP_NOZORDER   = 0x0004
	SWP_NOACTIVATE = 0x0010
)
