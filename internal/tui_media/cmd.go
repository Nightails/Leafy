package tui_media

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/nightails/leafy/internal/media"
	app "github.com/nightails/leafy/internal/tui_app"
)

type mediaMsg []string

func scanMediaCmd(paths []string) tea.Cmd {
	return func() tea.Msg {
		if len(paths) == 0 {
			return mediaMsg{}
		}
		m, err := media.GetMediaFiles(paths)
		if err != nil {
			return app.ErrMsg(err)
		}
		return mediaMsg(m)
	}
}
