package main

import (
	"fmt"
	render "hello/render"
	state "hello/state"
	utils "hello/utils"
	"log"
	"os"

	ui "github.com/gizak/termui/v3"
)

func main() {
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}

	err, branches := utils.GetGitBranches()
	if err != nil {
		ui.Close()
		fmt.Println("Error: No git repository found")
		os.Exit(1)
		return
	}
	_, currentBranch := utils.GetCurrentBranch()

	termWidth, termHeight := ui.TerminalDimensions()

	// Initial State
	s := state.InitState(
		branches,
		currentBranch,
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
		default:
			state.Multiplexer(e.ID, s, exit)
		}
	}

}
