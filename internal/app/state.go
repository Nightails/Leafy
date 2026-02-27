package app

type state struct {
	devices []device
	media   []medium
	tasks   []task
}

type device struct {
	name       string
	path       string
	mountpoint string
}

type medium struct {
	name          string
	format        string
	src, dest     string
	copied, total int64
}

type task struct {
	id          int
	media       medium
	done, total int64
	err         error
}
