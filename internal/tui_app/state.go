package tui_app

type AppState struct {
	MountPoints     []string
	FilesToTransfer []string
}

type AppStateMsg struct {
	State AppState
}

type DeviceMountedMsg struct {
	MountPoint string
}

type FileSelectedMsg struct {
	Path string
}
