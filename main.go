package main

import (
	"embed"
	"log"
	"os/exec"

	"xorriso-ui/pkg/xorriso"
	"xorriso-ui/services"

	"github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	xorrisoPath, err := exec.LookPath("xorriso")
	if err != nil {
		log.Fatal("xorriso not found in PATH. Please install xorriso (version 1.5.6+): sudo apt install xorriso / sudo zypper install xorriso")
	}

	executor := xorriso.NewExecutor(xorrisoPath)

	app := application.New(application.Options{
		Name:        "xorriso-ui",
		Description: "Modern disc burning GUI",
		Services: []application.Service{
			application.NewService(services.NewDeviceService(executor)),
			application.NewService(services.NewProjectService()),
			application.NewService(services.NewBurnService(executor)),
			application.NewService(services.NewSettingsService(executor)),
		},
		Assets: application.AssetOptions{
			Handler: application.AssetFileServerFS(assets),
		},
		Linux: application.LinuxOptions{
			ProgramName: "xorriso-ui",
		},
	})

	app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:            "xorriso-ui",
		Width:            1280,
		Height:           800,
		MinWidth:         900,
		MinHeight:        600,
		BackgroundColour: application.NewRGB(17, 24, 39),
		URL:              "/",
	})

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
