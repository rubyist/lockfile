package lockfile

type Locker interface {
	LockRead() error
	LockWrite() error
	LockReadB() error
	LockWriteB() error
	Unlock()
	Owner() int
	Remove()
}
