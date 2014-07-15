// +build linux darwin freebsd openbsd netbsd dragonfly
// +build go1.3

package lockfile

import (
	"os"
)

func NewLockfile(directory, name string) (Locker, error) {
	return NewFcntlLockfile(directory, name)
}

func NewLockfileFromFile(file *os.File) (Locker, error) {
	return NewFcntlLockfileFromFile(file)
}
