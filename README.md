timecat
-------

Utility for keeping time.

building
--------
```
go get ./...
go build
```

usage
-----
```
# Running as a timer
./timecat
# ctrl-c to stop

# A running timer on a command:
./timecat sleep 1

# Prepend timestamps while running a command:
./timecat --timestamp sudo apt-get update
```
