package main

import (
	"embed"
	"syscall"

	"mira-valorant/internal/config"
	"mira-valorant/internal/ui"
	"mira-valorant/internal/wailsapp"

	"github.com/lxn/win"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
)

//go:embed all:frontend/dist
var assets embed.FS

var (
	kernel32             = syscall.NewLazyDLL("kernel32.dll")
	procGetConsoleWindow = kernel32.NewProc("GetConsoleWindow")
)

func hideConsoleWindow() {
	hwnd, _, _ := procGetConsoleWindow.Call()
	if hwnd != 0 {
		win.ShowWindow(win.HWND(hwnd), win.SW_HIDE)
	}
}

func main() {
	hideConsoleWindow()

	state := config.NewAppState()
	state.SetDefaults()

	go ui.RunScopeWindow(state)

	app := wailsapp.NewApp(state)

	err := wails.Run(&options.App{
		Title:             "Mira Controller",
		Width:             420,
		Height:            520,
		DisableResize:     false,
		Assets:            assets,
		BackgroundColour:  &options.RGBA{R: 18, G: 18, B: 20, A: 255},
		Frameless:         false,
		AlwaysOnTop:       false,
		HideWindowOnClose: false,
		OnStartup:         func(ctx *wails.Context) {},
		Bind:              []interface{}{app},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
