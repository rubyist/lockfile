package lockfile

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
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

func (l *FLockfile) Lock() error {
	f, err := os.OpenFile(l.Path, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		return err
	}
	l.file = f

	err = syscall.Flock(int(f.Fd()), syscall.LOCK_EX|syscall.LOCK_NB)
	if err != nil {
		f.Close()
		return fmt.Errorf("Could not obtain file lock, owned by %d", l.Owner())
	}

	l.lockObtained = true
	f.Write([]byte(fmt.Sprintf("%d\n", os.Getpid())))
	f.Sync()

	return nil
}

func (l *FLockfile) Unlock() {
	if !l.lockObtained {
		return
	}

	syscall.Flock(int(l.file.Fd()), syscall.LOCK_UN)
	l.file.Close()
	os.Remove(l.Path)
}

func (l *FLockfile) Owner() int {
	data, err := ioutil.ReadFile(l.Path)
	if err != nil {
		return -1
	}

	pid, err := strconv.ParseInt(string(data), 0, 32)
	if err != nil {
		return -1
	}

	return int(pid)
}
