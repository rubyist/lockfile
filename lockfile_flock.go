package lockfile

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"syscall"
)

type FLockfile struct {
	Path         string
	file         *os.File
	lockObtained bool
}

func NewFLockfile(path string) *FLockfile {
	return &FLockfile{Path: path}
}

func NewFLockfileFromFile(file *os.File) *FLockfile {
	return &FLockfile{file: file}
}

func (l *FLockfile) LockRead() error {
	return l.lock(false, false)
}

func (l *FLockfile) LockWrite() error {
	return l.lock(true, false)
}

func (l *FLockfile) LockReadB() error {
	return l.lock(false, true)
}

func (l *FLockfile) LockWriteB() error {
	return l.lock(true, true)
}

func (l *FLockfile) Unlock() {
	if !l.lockObtained {
		return
	}

	syscall.Flock(int(l.file.Fd()), syscall.LOCK_UN)
	l.file.Close()
}

func (l *FLockfile) Owner() int {
	data, err := ioutil.ReadFile(l.Path)
	if err != nil {
		return -1
	}

	stripped := strings.Trim(string(data), "\n")
	pid, err := strconv.ParseInt(stripped, 0, 32)
	if err != nil {
		return -1
	}

	return int(pid)
}

func (l *FLockfile) lock(exclusive, blocking bool) error {
	if l.file == nil {
		f, err := os.OpenFile(l.Path, os.O_CREATE|os.O_RDWR, 0666)
		if err != nil {
			return err
		}
		l.file = f
	}

	var flags int
	if exclusive {
		flags = syscall.LOCK_EX
	} else {
		flags = syscall.LOCK_SH
	}
	if !blocking {
		flags |= syscall.LOCK_NB
	}

	err := syscall.Flock(int(l.file.Fd()), flags)
	if err != nil {
		l.file.Close()
		return ErrFailedToLock
	}

	l.lockObtained = true
	l.file.Write([]byte(fmt.Sprintf("%d\n", os.Getpid())))
	l.file.Sync()

	return nil
}
