package media

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/nightails/leafy/internal/app"
)

type mediaMsg []string

func scanMediaCmd(paths []string) tea.Cmd {
	return func() tea.Msg {
		if len(paths) == 0 {
			return mediaMsg{}
		}
		m, err := GetMediaFiles(paths)
		if err != nil {
			return app.ErrMsg(err)
		}
		return mediaMsg(m)
	}
}

type transferProgressMsg struct {
	Copied int64
	Total  int64
}

type transferDoneMsg struct{}

type transferStartedMsg struct {
	Ch <-chan tea.Msg
}

// transferMediaCmd starts copying in a goroutine and returns a message containing
// a channel that will emit progress/done/err messages.
func transferMediaCmd(src, dst string) tea.Cmd {
	return func() tea.Msg {
		ch := make(chan tea.Msg, 32)

		go func() {
			defer close(ch)

			send := func(msg tea.Msg) {
				select {
				case ch <- msg:
				default:
				}
			}

			err := copyFileWithProgress(src, dst, func(copied, total int64) {
				send(transferProgressMsg{Copied: copied, Total: total})
			})
			if err != nil {
				send(app.ErrMsg(fmt.Errorf("transfer: %w", err)))
				return
			}
			send(transferDoneMsg{})
		}()

		return transferStartedMsg{Ch: ch}
	}
}

// listenTransferCmd waits for the next message from the transfer goroutine.
func listenTransferCmd(ch <-chan tea.Msg) tea.Cmd {
	return func() tea.Msg {
		msg, ok := <-ch
		if !ok {
			return transferDoneMsg{}
		}
		return msg
	}
}
