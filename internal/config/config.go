package config

import (
	"image/color"
	"sync"
)

type AppState struct {
	sync.RWMutex
	isActive        bool
	zoomLevel       float64
	scopeSize       float64
	borderEnabled   bool
	borderColor     int
	borderThickness float64
	followCursor    bool
	refreshRate     int
	hotkeyMode      string // "Toggle" or "Hold"
	hotkey          int    // Virtual Key Code
}

func NewAppState() *AppState {
	return &AppState{
		refreshRate: 30,
	}
}

func (s *AppState) SetDefaults() {
	s.Lock()
	defer s.Unlock()
	s.isActive = false
	s.zoomLevel = 2.0
	s.scopeSize = 250.0
	s.borderEnabled = true
	s.borderColor = 0
	s.borderThickness = 2.0
	s.followCursor = false
	s.hotkeyMode = "Toggle"
	s.hotkey = 0 // None
}

// Getters
func (s *AppState) HotkeyMode() string {
	s.RLock()
	defer s.RUnlock()
	return s.hotkeyMode
}

func (s *AppState) Hotkey() int {
	s.RLock()
	defer s.RUnlock()
	return s.hotkey
}

func (s *AppState) IsActive() bool {
	s.RLock()
	defer s.RUnlock()
	return s.isActive
}

func (s *AppState) ZoomLevel() float64 {
	s.RLock()
	defer s.RUnlock()
	return s.zoomLevel
}

func (s *AppState) ScopeSize() float64 {
	s.RLock()
	defer s.RUnlock()
	return s.scopeSize
}

func (s *AppState) BorderEnabled() bool {
	s.RLock()
	defer s.RUnlock()
	return s.borderEnabled
}

func (s *AppState) BorderColor() int {
	s.RLock()
	defer s.RUnlock()
	return s.borderColor
}

func (s *AppState) BorderThickness() float64 {
	s.RLock()
	defer s.RUnlock()
	return s.borderThickness
}

func (s *AppState) FollowCursor() bool {
	s.RLock()
	defer s.RUnlock()
	return s.followCursor
}

func (s *AppState) RefreshRate() int {
	s.RLock()
	defer s.RUnlock()
	return s.refreshRate
}

// Setters
func (s *AppState) SetIsActive(v bool) {
	s.Lock()
	defer s.Unlock()
	s.isActive = v
}

func (s *AppState) SetZoomLevel(v float64) {
	s.Lock()
	defer s.Unlock()
	s.zoomLevel = v
}

func (s *AppState) SetScopeSize(v float64) {
	s.Lock()
	defer s.Unlock()
	s.scopeSize = v
}

func (s *AppState) SetBorderEnabled(v bool) {
	s.Lock()
	defer s.Unlock()
	s.borderEnabled = v
}

func (s *AppState) SetBorderColor(v int) {
	s.Lock()
	defer s.Unlock()
	s.borderColor = v
}

func (s *AppState) SetBorderThickness(v float64) {
	s.Lock()
	defer s.Unlock()
	s.borderThickness = v
}

func (s *AppState) SetFollowCursor(v bool) {
	s.Lock()
	defer s.Unlock()
	s.followCursor = v
}

func (s *AppState) SetHotkeyMode(v string) {
	s.Lock()
	defer s.Unlock()
	s.hotkeyMode = v
}

func (s *AppState) SetHotkey(v int) {
	s.Lock()
	defer s.Unlock()
	s.hotkey = v
}

func GetColorFromInt(c int) color.Color {
	switch c {
	case 0:
		return color.RGBA{255, 0, 0, 255} // Red
	case 1:
		return color.RGBA{0, 255, 0, 255} // Green
	case 2:
		return color.RGBA{0, 0, 255, 255} // Blue
	case 3:
		return color.RGBA{255, 255, 0, 255} // Yellow
	default:
		return color.RGBA{255, 255, 255, 255}
	}
}
