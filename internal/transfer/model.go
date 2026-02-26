package transfer

import (
	tea "github.com/charmbracelet/bubbletea"
)

type state int

const (
	queued = iota
	running
	done
	failed
	cancelled
)

type task struct {
	id          int
	src, dest   string
	total, done int64
	state       state
	err         error
}
type Model struct {
	tasks         []task
	maxConcurrent int
	runningCount  int
	indexes       []int
	doneCount     int
}

func NewModel() Model {
	return Model{
		tasks:         make([]task, 0),
		maxConcurrent: 4, // allowing 4 concurrent transfers
		runningCount:  0,
		indexes:       make([]int, 0),
		doneCount:     0,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m Model) View() string {
	return "Transferring"
}
