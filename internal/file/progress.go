package file

import (
	"io"
	"sync/atomic"
	"time"
)

type ProgressFn func(copied, total int64)

type progressWriter struct {
	w      io.Writer
	copied *atomic.Int64
}

// Write implements io.Writer.
func (p *progressWriter) Write(b []byte) (int, error) {
	n, err := p.w.Write(b)
	p.copied.Add(int64(n))
	return n, err
}

func startProgressWriter(copied *atomic.Int64, total int64, progress ProgressFn, interval time.Duration) func() {
	if progress == nil {
		return func() {}
	}

	done := make(chan struct{})
	finished := make(chan struct{})

	go func() {
		defer close(finished)

		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-done:
				progress(copied.Load(), total)
				return
			case <-ticker.C:
				progress(copied.Load(), total)
			}
		}
	}()

	return func() {
		close(done)
		<-finished
	}
}
