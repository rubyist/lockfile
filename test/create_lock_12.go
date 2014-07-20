package main

import (
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

	return lockfile.NewFLockfile(name), nil
}
