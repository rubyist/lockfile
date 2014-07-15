# lockfile

Lockfile provides a simple wrapper for `fcntl` and `flock` based advisory file locking.

## Installation

```
go get github.com/rubyist/lockfile
```

## Examples

`lockfile` provides a function `NewLockfile()` which will use an `fcntl` based lock, if
it is available on the system. `fcntl` based locks are available in Go 1.3 and later. If
`fcntl` is not available it will fall back to an `flock` based lock.

```go
lock := lockfile.NewLockfile("myfile")
```

A `fcntl` based lock can be built explicitly:

```go
lock := lockfile.NewFcntlLockfile("myfile")
```

A `flock` based lock can be built explicitly:

```go
lock := lockfile.NewFLockfile("myfile")
```

An `os.File` can also be used.

```go
file, _ := os.Open("myfile")
lock := lockfile.NewLockfileFromFile(file)
lock2 := lockfile.NewFcntlLockfileFromFile(file)
lock3 := lockfile.NewFLockfileFromFile(file)
```

### Read Locks

The `Read()` function will lock a file for reading. When a file is locked for reading,
other processes will be able to obtain read locks on the file. Write locks cannot be
obtained on files locked for reading. If the lock cannot be obtained, `Read()` will
return an error immediately

The `ReadB()` function provides the same read locking functionality but will block
until the lock can be obtained.

```go
err := lock.Read()
if err != nil {
  // Lock not obtained
}

lock.ReadB() // blocks until lock can be obtained
```

### Write Locks

The `Write()` function will lock a file for writing. When a file is locked for writing,
other processes will not be able to obtain read or write locks on the file. If the lock
cannot be obtained, `Write()` will return an error immediately.

The `WriteB()` function provides the same write locking functionality but will block
until the lock can be obtained.

```go
err := lock.Write()
if err != nil {
  // Lock not obtained
}

lock.WriteB() // blocks until the lock can be obtained

### Releasing locks

The `Unlock()` function will release a lock. Closing a file or exiting the process
will also release any locks held on the file.

## Bugs and Issues

Issues and pull requests on [GitHub](https://github.com/rubyist/lockfile)