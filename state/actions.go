package state

import (
	utils "hello/utils"
	"strings"
)

const FILTER_STATUS = "Enter: Confirm filter | Esc: Clear filter and quit filter mode"
const VIEW_STATUS = "d: Delete | n: New | /: Filter | c: Checkout | r: Rename | [q|esc|ctrl-c]: Quit"
const RENAME_STATUS = "Enter: Confirm rename | Esc: Exit rename mode"
const NEW_BRANCH_STATUS = "Enter: Confirm create branch| Esc: Exit new branch mode"

func Multiplexer(e string, s *State, exit func()) {
	switch m := s.GetMode(); m {
	case View:
		viewMode(e, s, exit)
	case Delete:
		deleteMode(e, s)
	case Filter:
		filterMode(e, s)
	case Rename:
		renameMode(e, s)
	case NewBranch:
		newMode(e, s)
	case Error:
		errorMode(e, s)
	default:
		// fmt.Println("default", m, e)
	}

	// Emit state change
	s.emit()
}

func Resize(w int, h int, s *State) {
	s.SetDimensions(w, h)
	s.emit()
}

// Modes
func viewMode(e string, s *State, exit func()) {
	i := s.GetSelectedRow()
	n := len(s.GetFilteredBranches())
	switch e {
	case "j", "<Down>":
		if i < n-1 {
			s.SetSelectedRow(i + 1)
		}
	case "k", "<Up>":
		if i > 0 {
			s.SetSelectedRow(i - 1)
		}
	case "<Home>":
		s.selectedRow = 0
	case "/":
		transitionToFilterMode(s)
	case "d":
		transitionToDeleteMode(s)
	case "r":
		transitionToRenameMode(s)
	case "n":
		transitionToNewBranchMode(s)
	case "c":
		checkout(s)
	case "G", "<End>":
		s.selectedRow = n - 1
	case "q", "<Escape>":
		exit()
	}
}

func filterMode(e string, s *State) {
	// Accept any git branch valid string
	f := s.GetFilter()
	if len(e) == 1 {
		newFilter := f + e
		setFilter(newFilter, s)
		return
	}

	switch e {
	case "<Backspace>":
		if len(f) > 0 {
			setFilter(f[:len(f)-1], s)
		}
	case "<Escape>":
		setFilter("", s)
		transitionToViewMode(s)
	case "<Enter>":
		transitionToViewMode(s)
	}
}

func renameMode(e string, s *State) {
	// Accept any git branch valid string
	rename := s.GetRename()
	if len(e) == 1 {
		setRename(rename+e, s)
		return
	}

	switch e {
	case "<Backspace>":
		if len(rename) > 0 {
			setRename(rename[:len(rename)-1], s)
		}
	case "<Escape>":
		setRename("", s)
		transitionToViewMode(s)
	case "<Enter>":
		// Execute rename
		name := s.GetFilteredBranches()[s.GetSelectedRow()]
		if err, out := utils.RenameGitBranch(name, rename); err != nil {
			transitionToErrorMode(out, s)
		} else {
			refreshBranches(s)
			transitionToViewMode(s)
		}
	}
}

func deleteMode(e string, s *State) {
	// Do a check to see if branch can be deleted here
	switch e {
	case "y", "Y":
		b := s.GetFilteredBranches()[s.GetSelectedRow()]
		if err, out := utils.DeleteGitBranch(b); err != nil {
			transitionToErrorMode(out, s)

		} else {
			transitionToViewMode(s)
		}
		refreshBranches(s)
	case "n", "N", "<Escape>":
		transitionToViewMode(s)
	}
}

func newMode(e string, s *State) {
	// Accept any git branch valid string
	newBranch := s.GetNewBranch()
	if len(e) == 1 {
		setNewBranch(newBranch+e, s)
		return
	}

	switch e {
	case "<Backspace>":
		if len(newBranch) > 0 {
			setNewBranch(newBranch[:len(newBranch)-1], s)
		}
	case "<Escape>":
		setNewBranch("", s)
		transitionToViewMode(s)
	case "<Enter>":
		// Execute rename
		base := s.GetFilteredBranches()[s.GetSelectedRow()]
		if err, out := utils.CreateGitBranch(newBranch, base); err != nil {
			transitionToErrorMode(out, s)
		} else {
			refreshBranches(s)
			transitionToViewMode(s)
		}
	}
}

func errorMode(e string, s *State) {
	switch e {
	case "q", "<Escape>":
		transitionToViewMode(s)
	}
}

func checkout(s *State) {
	if bs, r := s.GetFilteredBranches(), s.GetSelectedRow(); len(bs) > r {
		if err, out := utils.CheckoutGitBranch(bs[r]); err != nil {
			transitionToErrorMode(out, s)
		} else {
			s.SetCurrentBranch(bs[r])
		}
	}
}

// Transition
func transitionToViewMode(s *State) {
	s.SetMode(View)
	setViewStatus(s)
}

func transitionToFilterMode(s *State) {
	s.SetMode(Filter)
	setFilter("", s)
}

func transitionToDeleteMode(s *State) {
	i := s.GetSelectedRow()
	b := s.GetFilteredBranches()
	c := s.GetCurrentBranch()
	if i >= len(b) {
		transitionToErrorMode("Selection out of bounds", s)
	} else if b[i] == c {
		transitionToErrorMode("Unable to delete current branch", s)
	} else {
		s.SetStatus("Delete branch '" + b[i] + "'? (y,n)")
		s.SetMode(Delete)
	}
}

func transitionToRenameMode(s *State) {
	setRename("", s)
	s.SetMode(Rename)
}

func transitionToNewBranchMode(s *State) {
	setNewBranch("", s)
	s.SetMode(NewBranch)
}

func transitionToErrorMode(e string, s *State) {
	s.SetMode(Error)
	s.SetStatus(strings.TrimSpace(e) + "\nPress 'q' or <Escape> to go back to view mode")
}

// Helper functions -- prefix functions with 'set' to indicate side effects
func filterBranches(f string, branches []string) []string {
	fb := []string{}
	for _, b := range branches {
		if strings.Contains(b, f) {
			fb = append(fb, b)
		}
	}
	return fb
}

func setFilter(f string, s *State) {
	_, b := utils.GetGitBranches()
	s.SetFilteredBranches(filterBranches(f, b))
	s.SetFilter(f)
	s.SetStatus("Filter: " + f + "█" + "\n" + FILTER_STATUS)
}

func setRename(r string, s *State) {
	n := s.GetFilteredBranches()[s.GetSelectedRow()]
	s.SetRename(r)
	s.SetStatus("Rename: " + n + " → " + r + "█" + "\n" + RENAME_STATUS)
}

func setNewBranch(newBranch string, s *State) {
	base := s.GetFilteredBranches()[s.GetSelectedRow()]
	s.SetNewBranch(newBranch)
	s.SetStatus("Base: " + base + " → " + newBranch + "█" + "\n" + NEW_BRANCH_STATUS)
}

func setViewStatus(s *State) {
	f := s.GetFilter()
	if f != "" {
		s.SetStatus("Filter: " + f + "\n" + VIEW_STATUS)
	} else {
		s.SetStatus(VIEW_STATUS)
	}
}

func refreshBranches(s *State) {
	_, branches := utils.GetGitBranches()
	_, currentBranch := utils.GetCurrentBranch()
	s.SetBranches(branches)
	s.SetFilteredBranches(filterBranches(s.GetFilter(), branches))
	s.SetCurrentBranch(currentBranch)
}
