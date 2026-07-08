package watcher

import "golang.org/x/sys/unix"

func addWatch(
	fileDispector int,
	path string,
	mask uint32,
) (int, error) {
	return unix.InotifyAddWatch(
		fileDispector,
		path,
		mask,
	)
}
