package media

import (
	"fmt"
	"io"
	"os"
	"time"

	"golang.org/x/sys/unix"
)

func copyFileWithProgress(src, dst string, progress progressFn) error {
	// reserve permission
	si, err := os.Lstat(src)
	if err != nil {
		return fmt.Errorf("stat src: %w", err)
	}
	if !si.Mode().IsRegular() {
		return fmt.Errorf("src is not a regular file: %s (%v)", src, si.Mode())
	}
	total := si.Size()

	// reserve timestamps
	atime, mtime, err := statATimeMTime(src)
	if err != nil {
		return fmt.Errorf("get atime and mtime: %w", err)
	}

	if err := os.MkdirAll(dst, 0755); err != nil {
		return fmt.Errorf("mkdir dst dir: %w", err)
	}

	in, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("open src: %w", err)
	}
	defer func(in *os.File) {
		_ = in.Close()
	}(in)

	// write to a temporary file
	tmp := dst + ".tmp"
	out, err := os.OpenFile(tmp, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, si.Mode().Perm())
	if err != nil {
		return fmt.Errorf("create dst: %w", err)
	}

	pw := newProgressWriter(out, total, progress, 150*time.Millisecond)
	_, copyErr := io.Copy(pw, in)

	syncErr := out.Sync()
	closeErr := out.Close()

	if copyErr != nil {
		_ = os.Remove(tmp)
		return fmt.Errorf("copy: %w", copyErr)
	}
	if syncErr != nil {
		_ = os.Remove(tmp)
		return fmt.Errorf("sync: %w", syncErr)
	}
	if closeErr != nil {
		_ = os.Remove(tmp)
		return fmt.Errorf("close tmp: %w", closeErr)
	}

	pw.finish()

	if err := os.Rename(tmp, dst); err != nil {
		_ = os.Remove(dst)
		return fmt.Errorf("rename tmp to dst: %w", err)
	}
	if err := os.Chmod(dst, si.Mode().Perm()); err != nil {
		return fmt.Errorf("chmod dst: %w", err)
	}
	if err := setATimeMTime(dst, atime, mtime); err != nil {
		return fmt.Errorf("set dst timestamps: %w", err)
	}

	return nil
}

func statATimeMTime(path string) (atime, mtime unix.Timespec, err error) {
	var st unix.Stat_t
	if err := unix.Stat(path, &st); err != nil {
		return unix.Timespec{}, unix.Timespec{}, err
	}
	return st.Atim, st.Mtim, nil
}

func setATimeMTime(path string, atime, mtime unix.Timespec) error {
	return unix.UtimesNano(path, []unix.Timespec{atime, mtime})
}
