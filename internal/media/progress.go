package media

import (
	"io"
	"time"
)

type progressFn func(copied, total int64)

type progressWriter struct {
	w        io.Writer
	total    int64
	progress progressFn

	interval time.Duration
	lastEmit time.Time
	copied   int64
	enabled  bool
}

func newProgressWriter(w io.Writer, total int64, progress progressFn, interval time.Duration) *progressWriter {
	return &progressWriter{
		w:        w,
		total:    total,
		progress: progress,
		interval: interval,
		lastEmit: time.Now(),
		enabled:  progress != nil,
	}
}

func (p *progressWriter) Write(b []byte) (int, error) {
	n, err := p.w.Write(b)
	p.copied += int64(n)
	if p.enabled && time.Since(p.lastEmit) >= p.interval {
		p.progress(p.copied, p.total)
		p.lastEmit = time.Now()
	}
	return n, err
}

func (p *progressWriter) finish() {
	if p.enabled {
		p.progress(p.copied, p.total)
	}
}
