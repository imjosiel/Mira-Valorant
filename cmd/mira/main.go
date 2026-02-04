package main

import (
	"mira-valorant/internal/config"
	"mira-valorant/internal/ui"
)

func main() {
	// 1. Initialize State
	state := config.NewAppState()
	state.SetDefaults()

	// 2. Start Scope Window in a separate goroutine
	// The scope window needs its own thread for the message loop if we want it to be independent,
	// but UI threads usually need to be the main thread or carefully managed.
	// Since we have two windows (Control and Scope), we need two loops or one managed loop.
	// lxn/walk runs its own loop.
	// We'll run the Scope logic in a goroutine that locks its own thread.
	go ui.RunScopeWindow(state)

	// 3. Start Control Window (Blocks Main Thread)
	ui.RunControlWindow(state)
}
