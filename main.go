package main

import (
	render "hello/render"
	state "hello/state"
	utils "hello/utils"
	"log"

	ui "github.com/gizak/termui/v3"
)

func main() {
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}

	// defer ui.Close()

	termWidth, termHeight := ui.TerminalDimensions()

	// Initial State
	s := state.InitState(
		utils.GetGitBranches(),
		utils.GetCurrentBranch(),
		0,
		state.VIEW_STATUS,
		state.TerminalDimensions{
			H: termHeight, W: termWidth,
		},
	)

	// Subscribe to state changes
	// state will call render when there is a change in state
	render.Render(s)
	s.Subscribe(render.Render)

	// Main loop
	uiEvents := ui.PollEvents()

	// Exit callback
	shouldExit := false
	exit := func() {
		shouldExit = true
		ui.Close()
	}

	for {
		if shouldExit {
			return
		}
		e := <-uiEvents
		switch e.ID {
		// forced exit
		case "<C-c>":
			exit()
		case "<Resize>":
			payload := e.Payload.(ui.Resize)
			state.Resize(payload.Width, payload.Height, s)
			// grid.SetRect(0, 0, payload.Width, payload.Height-STATUS_HEIGHT)
			// p.SetRect(0, payload.Height-STATUS_HEIGHT, payload.Width, payload.Height)
			// ui.Clear()
			// ui.Render(grid, p)
		default:
			state.Multiplexer(e.ID, s, exit)
		}

		// Render
		// ui.Render(grid, p)
	}

}
