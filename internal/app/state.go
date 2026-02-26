package app

type State struct {
	MountPoints []string
	MediaFiles  []MediaFile
}

type MediaFile struct {
	Name string
	Src  string
	Dest string
}
