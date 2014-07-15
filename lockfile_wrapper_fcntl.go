// +build linux darwin freebsd openbsd netbsd dragonfly
// +build go1.3

package lockfile

func NewLockfile(directory, name string) (Locker, error) {
	return NewFcntlLockfile(directory, name)
}
