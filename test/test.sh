#!/bin/sh
rm -f test/test
cd test
if [[ `go version` == *go1.3* ]]
then go build test.go create_lock_13.go
else go build test.go create_lock_12.go
fi

function allows() {
  printf "$1 $2 allows $3 $4..."
  ./test -type $1 -access $2 -wait &
  p=$!
  sleep .1
  ./test -type $3 -access $4
  if [ $? -eq 0 ]
  then printf "OK\n"
  else printf "FAILED\n"
  fi
  kill $p
}

function denies() {
  printf "$1 $2 denies $3 $4..."
  ./test -type $1 -access $2 -wait &
  p=$!
  sleep .1
  ./test -type $3 -access $4
  if [ $? -eq 2 ]
  then printf "OK\n"
  else printf "FAILED\n"
  fi
  kill $p
}

allows flock read flock read
denies flock read flock write
denies flock write flock read
denies flock write flock write

if [[ `go version` == *go1.3* ]]
then
  allows fcntl read fcntl read
  denies fcntl read fcntl write
  denies fcntl write fcntl read
  denies fcntl write fcntl write

  allows flock read fcntl read
  denies flock read fcntl write
  denies flock write fcntl read
  denies flock write fcntl write

  allows fcntl read flock read
  denies fcntl read flock write
  denies fcntl write flock read
  denies fcntl write flock write
fi

rm -f lockfiletest.lock
rm -f test
