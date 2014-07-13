package lockfile

type Locker interface {
	Lock() error
	Unlock()
	Owner() int
}
