// +build go1.2,!go1.3

package lockfile

func NewLockfile(directory, name string) (Locker, error) {
	return NewFLockfile(directory, name)
}
