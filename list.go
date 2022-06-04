//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"log"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

var previousKey string

func visualHandler(l *widgets.List, eid string) {
	switch eid {
	case "q", "<C-c>":
		return
	case "j", "<Down>":
		l.ScrollDown()
	case "k", "<Up>":
		l.ScrollUp()
	case "<C-d>":
		l.ScrollHalfPageDown()
	case "<C-u>":
		l.ScrollHalfPageUp()
	case "<C-f>":
		l.ScrollPageDown()
	case "<C-b>":
		l.ScrollPageUp()
	case "g":
		if previousKey == "g" {
			l.ScrollTop()
		}
	case "<Home>":
		l.ScrollTop()
	case "G", "<End>":
		l.ScrollBottom()
	}

	if previousKey == "g" {
		previousKey = ""
	} else {
		previousKey = eid
	}
	fmt.Println(previousKey)
}

func main() {
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	l := widgets.NewList()
	l.Title = "List"
	l.Rows = []string{
		"[0] github.com/gizak/termui/v3",
		"[1] [你好，世界](fg:blue)",
		"[2] [こんにちは世界](fg:red)",
		"[3] [color](fg:white,bg:green) output",
		"[4] output.go",
		"[5] random_out.go",
		"[6] dashboard.go",
		"[7] foo",
		"[8] bar",
		"[9] baz",
	}
	l.TextStyle = ui.NewStyle(ui.ColorYellow)
	l.WrapText = false
	l.SetRect(0, 0, 25, 5)

	ui.Render(l)

	uiEvents := ui.PollEvents()
	for {
		e := <-uiEvents
		// List control handler
		visualHandler(l, e.ID)

		ui.Render(l)
	}
}
