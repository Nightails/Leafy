package app

import tea "github.com/charmbracelet/bubbletea"

type (
	errMsg     error
	devicesMsg []device
	mediaMsg   []medium
)

type copyStartedMsg struct {
	Ch <-chan tea.Msg
}

type copyProgressMsg struct {
	Index  int
	Copied int64
	Total  int64
}

type copyDoneMsg struct {
	Index int
}

type copyFinishedMsg struct{}
type copyErrorMsg struct {
	Index int
	Err   error
}

type deleteDoneMsg struct{}
