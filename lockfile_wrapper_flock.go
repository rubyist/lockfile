// +build go1.2,!go1.3

package lockfile

import (
	"os"
)

func NewLockfile(directory, name string) (Locker, error) {
	return NewFLockfile(directory, name)
}

func NewLockfileFromFile(file *os.File) (Locker, error) {
	return NewFLockfileFromFile(file)
}
