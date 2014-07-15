// +build linux darwin freebsd openbsd netbsd dragonfly
// +build go1.3

package lockfile

import (
	"fmt"
	"os"
	"path/filepath"
	"syscall"
)

type FcntlLockfile struct {
	Path         string
	file         *os.File
	lockObtained bool
	ft           *syscall.Flock_t
}

func NewFcntlLockfile(directory, name string) (*FcntlLockfile, error) {
	fileName := fmt.Sprintf("%s.lock", name)
	return &FcntlLockfile{Path: filepath.Join(directory, fileName)}, nil
}

func NewFcntlLockfileFromFile(file *os.File) (*FcntlLockfile, error) {
	return &FcntlLockfile{file: file}, nil
}

func (l *FcntlLockfile) LockRead() error {
	return l.lock(false, false)
}

func (l *FcntlLockfile) LockWrite() error {
	return l.lock(true, false)
}

func (l *FcntlLockfile) LockReadB() error {
	return l.lock(false, true)
}

func (l *FcntlLockfile) LockWriteB() error {
	return l.lock(true, true)
}

func (l *FcntlLockfile) Unlock() {
	if !l.lockObtained {
		return
	}

	l.ft.Type = syscall.F_UNLCK
	syscall.FcntlFlock(l.file.Fd(), syscall.F_SETLK, l.ft)
	l.file.Close()
}

func (l *FcntlLockfile) Owner() int {
	ft := &syscall.Flock_t{}
	*ft = *l.ft

	err := syscall.FcntlFlock(l.file.Fd(), syscall.F_GETLK, ft)
	if err != nil {
		fmt.Println(err)
		return -1
	}

	if ft.Type == syscall.F_UNLCK {
		fmt.Println(err)
		return -1
	}

	return int(ft.Pid)
}

func (l *FcntlLockfile) Remove() {
	os.Remove(l.Path)
}

func (l *FcntlLockfile) lock(exclusive, blocking bool) error {
	if l.lockObtained {
		return fmt.Errorf("Already locked")
	}

	if l.file == nil {
		f, err := os.OpenFile(l.Path, os.O_CREATE|os.O_RDWR, 0666)
		if err != nil {
			return err
		}
		l.file = f
	}

	ft := &syscall.Flock_t{
		Whence: int16(os.SEEK_SET),
		Start:  0,
		Len:    0,
		Pid:    int32(os.Getpid()),
	}
	l.ft = ft

	if exclusive {
		ft.Type = syscall.F_WRLCK
	} else {
		ft.Type = syscall.F_RDLCK
	}
	var flags int
	if blocking {
		flags = syscall.F_SETLKW
	} else {
		flags = syscall.F_SETLK
	}

	err := syscall.FcntlFlock(l.file.Fd(), flags, l.ft)
	if err != nil {
		owner := l.Owner()
		l.file.Close()
		return fmt.Errorf("Could not obtain file lock, owned by %d", owner)
	}

	l.lockObtained = true

	return nil
}
