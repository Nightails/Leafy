package transfer

import (
	tea "github.com/charmbracelet/bubbletea"
)

func startCopyCmd(task task) tea.Cmd {
	return func() tea.Msg {
		taskCh := make(chan tea.Msg)
		go func() {
			// TODO: Implement file transfer logic
			if task.err != nil {
				taskCh <- taskDoneMsg{ID: task.id}
				return
			}
			for task.done < task.total {
				taskCh <- taskProgressMsg{task.id, task.done, task.total}
			}
			taskCh <- taskDoneMsg{ID: task.id}
		}()
		return transferStartedMsg{Ch: taskCh}
	}
}

func listenTaskCmd(task task, ch chan tea.Msg) tea.Cmd {
	return func() tea.Msg {
		msg, ok := <-ch
		if !ok {
			return taskDoneMsg{ID: task.id}
		}
		return msg
	}
}
