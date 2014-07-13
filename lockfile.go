package lockfile

type Locker interface {
	LockRead() error
	LockWrite() error
	Unlock()
	Owner() int
	Remove()
}
