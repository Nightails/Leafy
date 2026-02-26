package transfer

import tea "github.com/charmbracelet/bubbletea"

type transferStartedMsg struct {
	Ch <-chan tea.Msg
}

type taskStartedMsg struct {
	ID int
}

type taskProgressMsg struct {
	ID    int
	Done  int64
	Total int64
}

type taskDoneMsg struct {
	ID int
}

type taskErrMsg struct {
	ID    int
	Error error
}

type taskCancelMsg struct {
	ID int
}
