package main

import (
	"fmt"
	"github.com/rubyist/lockfile"
	"os"
	"path/filepath"
)

func createLock(lockType string) (lockfile.Locker, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	name := filepath.Join(wd, "lockfiletest.lock")

	switch lockType {
	case "fcntl":
		return lockfile.NewFcntlLockfile(name), nil
	case "flock":
		return lockfile.NewFLockfile(name), nil
	}
	return nil, fmt.Errorf("invalid lock file type")
}
