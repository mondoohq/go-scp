//go:build openbsd || dragonfly || solaris

package scp

import (
	"os"
	"syscall"
	"time"
)

// These platforms spell the access time field Atim, like linux does.

func newFileInfoFromOS(fi os.FileInfo, replaceName string) *FileInfo {
	var name string
	if replaceName == "" {
		name = fi.Name()
	} else {
		name = replaceName
	}

	modTime := fi.ModTime()

	var accessTime time.Time
	sysStat, ok := fi.Sys().(*syscall.Stat_t)
	if ok {
		sec, nsec := sysStat.Atim.Unix()
		accessTime = time.Unix(sec, nsec)
	}

	return NewFileInfo(name, fi.Size(), fi.Mode(), modTime, accessTime)
}
