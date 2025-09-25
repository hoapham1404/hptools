package ui

import (
	"time"

	"github.com/wailsapp/wails/v3/pkg/application"

	"hptools/internal/config"
)

// SetupSystray initializes the system tray, menu and attaches window behavior.
// It returns a cleanup function that can be called on shutdown (currently no-op but left for future use).
func SetupSystray(app *application.App, win application.Window, cfg config.SystrayConfig) func() {
	systray := app.SystemTray.New()
	systray.SetLabel(cfg.Label)

	// Attach and configure window linking
	systray.AttachWindow(win)
	systray.WindowOffset(cfg.WindowOffset)
	systray.WindowDebounce(time.Duration(cfg.DebounceMS) * time.Millisecond)

	// Build menu
	menu := application.NewMenu()
	menu.Add("Open").OnClick(func(*application.Context) {
		win.Show()
	})

	menu.Add("Quit").OnClick(func(*application.Context) {
		app.Quit()
	})

	systray.SetMenu(menu)

	// Return a cleanup function placeholder
	return func() {
		// Future cleanup like removing icons, saving state, etc.
	}
}
