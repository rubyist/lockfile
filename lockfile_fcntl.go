// +build go1.3

package lockfile

import (
	"os"
	"syscall"
)

type FcntlLockfile struct {
	Path         string
	file         *os.File
	lockObtained bool
	ft           *syscall.Flock_t
}

func (f *FcntlLockfile) Lock() error {
	ft := &syscall.Flock_t{}
	ft.Start = 0
	ft.Len = 0
	ft.Pid = os.Getpid()
	ft.Whence = os.SEEK_SET

	f, err := os.OpenFile(l.Path, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		return err
	}
	l.file = f
	l.ft = ft

	err = syscall.FcntlFlock(f.Fd(), syscall.F_SETLK, ft)
	if err != nil {
		f.Close()
		return fmt.Errorf("Could not obtain file lock, owned by %s", l.Owner())
	}

	l.lockObtained = true

	return nil
}

func (f *FcntlLockfile) Unlock() {
	if !l.lockObtained {
		return
	}

	f.ft.Type = syscall.F_UNLK
	syscall.FcntlFlock(f.Fd(), syscall.F_SETLK, ft)
	f.file.Close()
}

func (f *FcntlLockfile) Owner() int {
	ft := &syscall.Flock_t{}
	*ft = *f.ft

	err = syscall.FcntlFlock(f.Fd(), syscall.F_GETLK, ft)
	if ft.Type == syscall.F_UNLK {
		return -1
	}

	return ft.Pid
}
