package state

type Mode int64

const (
	View Mode = iota
	Delete
	NewBranch
	Filter
	Checkout
	Rename
	Error
)

type TerminalDimensions struct {
	W int
	H int
}

type Subscriber func(s *State)

type State struct {
	branches         []string // true source
	filteredBranches []string // used for displaying
	currentBranch    string
	mode             Mode
	filter           string
	status           string
	rename           string
	newBranch        string
	selectedRow      int
	td               TerminalDimensions
	subscriber       Subscriber
}

var _state *State

// Initialization
func InitState(branches []string, currentBranch string, selectedRow int, status string, td TerminalDimensions) *State {
	s := State{
		branches,
		branches,
		currentBranch,
		0,
		"",
		status,
		"",
		"",
		selectedRow,
		td,
		nil,
	}
	_state = &s
	return _state
}

// Subscription

func (s *State) Subscribe(f Subscriber) {
	s.subscriber = f
}

func (s *State) emit() {
	s.subscriber(s)
}

// Getters and Setters
func (s *State) GetSelectedRow() int {
	return s.selectedRow
}

func (s *State) SetSelectedRow(row int) {
	s.selectedRow = row
}

func (s *State) GetBranches() []string {
	return s.branches
}

func (s *State) SetBranches(branches []string) {
	s.branches = branches
}

func (s *State) GetFilteredBranches() []string {
	return s.filteredBranches
}

func (s *State) SetFilteredBranches(branches []string) {
	s.filteredBranches = branches
}

func (s *State) GetDimensions() TerminalDimensions {
	return s.td
}

func (s *State) SetDimensions(width int, height int) {
	s.td = TerminalDimensions{width, height}
}

func (s *State) GetMode() Mode {
	return s.mode
}

func (s *State) SetMode(m Mode) {
	s.mode = m
}

func (s *State) GetFilter() string {
	return s.filter
}

func (s *State) SetFilter(f string) {
	s.filter = f
}

func (s *State) GetCurrentBranch() string {
	return s.currentBranch
}

func (s *State) SetCurrentBranch(b string) {
	s.currentBranch = b
}

func (s *State) GetStatus() string {
	return s.status
}

func (s *State) SetStatus(status string) {
	s.status = status
}

func (s *State) GetRename() string {
	return s.rename
}

func (s *State) SetRename(name string) {
	s.rename = name
}

func (s *State) GetNewBranch() string {
	return s.newBranch
}

func (s *State) SetNewBranch(branch string) {
	s.newBranch = branch
}

// Returns empty string when selected row is out of bounds
func (s *State) GetSelectedBranch() string {
	f := s.GetFilteredBranches()
	if l, i := len(f), s.GetSelectedRow(); l > i {
		return f[i]
	}
	return ""
}
