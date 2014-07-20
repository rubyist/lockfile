package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	lockTypePtr := flag.String("type", "flock", "lock type: <fcntl|flock>")
	accessPtr := flag.String("access", "read", "access: <read|write>")
	waitPtr := flag.Bool("wait", false, "wait until killed")

	flag.Parse()

	if !(*lockTypePtr == "flock" || *lockTypePtr == "fcntl") {
		os.Exit(1)
	}

	if !(*accessPtr == "read" || *accessPtr == "write") {
		os.Exit(1)
	}

	lock, err := createLock(*lockTypePtr)
	if err != nil {
		os.Exit(1)
	}

	if *accessPtr == "read" {
		err = lock.LockRead()
		if err != nil {
			os.Exit(2)
		}
	}
	if *accessPtr == "write" {
		err = lock.LockWrite()
		if err != nil {
			os.Exit(2)
		}
	}

	if *waitPtr {
		die := make(chan os.Signal, 1)
		signal.Notify(die, syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM)

		for {
			select {
			case <-die:
				lock.Unlock()
				os.Exit(0)
			}
			time.Sleep(time.Second)
		}
	} else {
		lock.Unlock()
		os.Exit(0)
	}

}
