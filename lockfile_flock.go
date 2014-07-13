package lockfile

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
)

type FLockfile struct {
	Path         string
	file         *os.File
	lockObtained bool
}

func NewFLockfile(directory, name string) (*FLockfile, error) {
	fileName := fmt.Sprintf("%s.lock", name)
	return &FLockfile{Path: filepath.Join(directory, fileName)}, nil
}

func (l *FLockfile) LockRead() error {
	return l.lock(syscall.LOCK_SH)
}

func (l *FLockfile) LockWrite() error {
	return l.lock(syscall.LOCK_EX)
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

func (l *FLockfile) Remove() {
	os.Remove(l.Path)
}

func (l *FLockfile) lock(lockType int) error {
	if l.lockObtained {
		return fmt.Errorf("Already locked")
	}

	f, err := os.OpenFile(l.Path, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		return err
	}
	l.file = f

	err = syscall.Flock(int(f.Fd()), lockType|syscall.LOCK_NB)
	if err != nil {
		f.Close()
		return fmt.Errorf("Could not obtain file lock, owned by %d", l.Owner())
	}

	l.lockObtained = true
	f.Write([]byte(fmt.Sprintf("%d\n", os.Getpid())))
	f.Sync()

	return nil
}
