package wailsapp

import (
	"mira-valorant/internal/config"
)

type App struct {
	state *config.AppState
}

func NewApp(state *config.AppState) *App {
	return &App{state: state}
}

type UIState struct {
	IsActive        bool
	ZoomLevel       float64
	ScopeSize       float64
	BorderEnabled   bool
	BorderColor     int
	BorderThickness float64
	FollowCursor    bool
	HotkeyMode      string
	Hotkey          int
	RenderBackend   string
}

func (a *App) GetState() UIState {
	return UIState{
		IsActive:        a.state.IsActive(),
		ZoomLevel:       a.state.ZoomLevel(),
		ScopeSize:       a.state.ScopeSize(),
		BorderEnabled:   a.state.BorderEnabled(),
		BorderColor:     a.state.BorderColor(),
		BorderThickness: a.state.BorderThickness(),
		FollowCursor:    a.state.FollowCursor(),
		HotkeyMode:      a.state.HotkeyMode(),
		Hotkey:          a.state.Hotkey(),
		RenderBackend:   a.state.RenderBackend(),
	}
}

func (a *App) SetActive(on bool) {
	a.state.SetIsActive(on)
}

func (a *App) ToggleActive() {
	a.state.SetIsActive(!a.state.IsActive())
}

func (a *App) SetZoomLevel(v float64) {
	a.state.SetZoomLevel(v)
}

func (a *App) SetBorderEnabled(v bool) {
	a.state.SetBorderEnabled(v)
}

func (a *App) SetBorderColor(v int) {
	a.state.SetBorderColor(v)
}

func (a *App) SetBorderThickness(v float64) {
	a.state.SetBorderThickness(v)
}

func (a *App) SetScopeSize(v float64) {
	a.state.SetScopeSize(v)
}

func (a *App) SetFollowCursor(v bool) {
	a.state.SetFollowCursor(v)
}

func (a *App) SetHotkeyMode(v string) {
	a.state.SetHotkeyMode(v)
}

func (a *App) SetHotkey(v int) {
	a.state.SetHotkey(v)
}

func (a *App) SetRenderBackend(v string) {
	a.state.SetRenderBackend(v)
}
