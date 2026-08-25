//go:build !linux && !darwin && !windows && !freebsd && !netbsd && !openbsd && !dragonfly && !solaris

package scp

import (
	"os"
	"time"
)

// Anything with no platform-specific access time accessor, such as aix.
// os.FileInfo exposes no portable access time, so it is left zero rather than
// failing to build; every other field is still reported.

func newFileInfoFromOS(fi os.FileInfo, replaceName string) *FileInfo {
	var name string
	if replaceName == "" {
		name = fi.Name()
	} else {
		name = replaceName
	}

	return NewFileInfo(name, fi.Size(), fi.Mode(), fi.ModTime(), time.Time{})
}
