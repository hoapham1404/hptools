package main

import (
	"embed"
	"log"

	"github.com/wailsapp/wails/v3/pkg/application"

	"hptools/internal/config"
	"hptools/internal/logging"
	"hptools/internal/services"
	"hptools/internal/windows"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// Load configuration
	cfg, err := config.Load(config.GetConfigPath())
	if err != nil {
		log.Printf("Warning: Failed to load config, using defaults: %v", err)
		cfg = config.Default()
	}

	// Setup logging
	logger := logging.NewLogger(cfg)
	appLogger := logging.WithComponent(logger, "app")

	appLogger.Info("Starting HP Tools", "version", "1.0.0")

	// Initialize Windows API
	api := windows.NewAPI()

	// Create services
	windowService := services.NewWindowService(api, logging.WithComponent(logger, "window_service"))
	wailsService := services.NewWailsWindowService(windowService)

	// Create Wails application
	app := application.New(application.Options{
		Name:        cfg.App.Name,
		Description: cfg.App.Description,
		Services: []application.Service{
			application.NewService(wailsService),
		},
		Assets: application.AssetOptions{
			Handler: application.AssetFileServerFS(assets),
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
	})

	// Create main window
	app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:     cfg.Window.Title,
		Width:     cfg.Window.Width,
		Height:    cfg.Window.Height,
		MinWidth:  cfg.Window.MinWidth,
		MinHeight: cfg.Window.MinHeight,
		MaxWidth:  cfg.Window.MaxWidth,
		Mac: application.MacWindow{
			InvisibleTitleBarHeight: cfg.Window.Mac.InvisibleTitleBarHeight,
			Backdrop:                cfg.Window.Mac.Backdrop,
			TitleBar:                cfg.Window.Mac.TitleBar,
		},
		BackgroundColour: cfg.Window.BackgroundColour,
		URL:              "/",
	})

	appLogger.Info("Application initialized, starting...")

	// Run the application
	if err := app.Run(); err != nil {
		appLogger.Error("Application failed", "error", err)
		log.Fatal(err)
	}

	appLogger.Info("Application stopped")
}
