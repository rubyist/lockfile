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

func (l *FcntlLockfile) LockRead() error {
	return l.lock(syscall.F_RDLCK)
}

func (l *FcntlLockfile) LockWrite() error {
	return l.lock(syscall.F_WRLCK)
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

func (l *FcntlLockfile) lock(lockType int16) error {
	if l.lockObtained {
		return fmt.Errorf("Already locked")
	}

	f, err := os.OpenFile(l.Path, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		return err
	}
	l.file = f

	ft := &syscall.Flock_t{
		Type:   lockType,
		Whence: int16(os.SEEK_SET),
		Start:  0,
		Len:    0,
		Pid:    int32(os.Getpid()),
	}
	l.ft = ft

	err = syscall.FcntlFlock(f.Fd(), syscall.F_SETLK, ft)
	if err != nil {
		owner := l.Owner()
		f.Close()
		return fmt.Errorf("Could not obtain file lock, owned by %d", owner)
	}

	l.lockObtained = true

	return nil
}
