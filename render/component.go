package render

import (
	"math"
	"strings"

	state "github.com/banghui/branch-manager/state"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

const STATUS_HEIGHT = 4

// Root Component
func Render(state *state.State) {
	// get the slice of components
	g := grid(state)
	sb := statusBar(state)

	// pass in states and all
	// get their ui.Drawable outputs
	// and render them
	ui.Render(g, sb)
}

func list(state *state.State) ui.Drawable {
	bs := state.GetFilteredBranches()
	cb := state.GetCurrentBranch()
	vbs := make([]string, len(bs))
	for i, b := range bs {
		if b == cb {
			vbs[i] = b + " ← current "
		} else {
			vbs[i] = b
		}
	}
	l := widgets.NewList()
	l.Rows = vbs
	l.Title = "Branches"
	l.SelectedRow = state.GetSelectedRow()
	l.TextStyle = ui.NewStyle(ui.ColorYellow)
	l.WrapText = false
	l.SelectedRowStyle.Bg = ui.ColorWhite
	l.SelectedRowStyle.Fg = ui.ColorBlack
	return l
}

func statusBar(s *state.State) ui.Drawable {
	sb := widgets.NewParagraph()
	sb.Title = "Status"
	sb.Border = true
	dim := s.GetDimensions()
	statusHeight := getStatusHeight(s)
	sb.SetRect(0, dim.H-statusHeight, dim.W, dim.H)
	sb.TextStyle.Fg = ui.ColorMagenta
	status := s.GetStatus()
	sb.Text = status
	return sb
}

func grid(s *state.State) ui.Drawable {
	l := list(s)
	grid := ui.NewGrid()
	dim := s.GetDimensions()
	statusHeight := getStatusHeight(s)
	grid.SetRect(0, 0, dim.W, dim.H-statusHeight)
	grid.Set(ui.NewRow(1.0, ui.NewCol(1.0, l)))
	return grid
}

// helper
func getStatusHeight(s *state.State) int {
	status := s.GetStatus()
	dim := s.GetDimensions()
	lines := strings.Split(status, "\n")
	rows := int(math.Max(2, float64(len(lines))))
	for _, s := range lines {
		if l := len(s); l+2 > dim.W {
			rows += 1
		}
	}
	return rows + 2
}
