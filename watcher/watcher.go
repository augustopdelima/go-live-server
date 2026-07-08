package watcher

import (
	"golang.org/x/sys/unix"
)

type AddWatchFunc func(
	fileDescriptor int,
	path string,
	mask uint32,
) (int, error)

type Watcher struct {
	fileDescriptor int
	root string
	watches map[int] string
	paths map[string]int
	events chan Event
	addWatch AddWatchFunc
}

func NewWatcher(root string) (*Watcher, error) {
	fileDescriptor, err := unix.InotifyInit1(unix.IN_CLOEXEC)

	if err != nil {
		return nil, err
	}

	watcher := &Watcher{
		fileDescriptor: fileDescriptor,
		root: root,
		watches: make(map[int]string),
		events: make(chan Event),
		paths: make(map[string]int),
		addWatch: addWatch,
	}

	return watcher, nil
}

const (
	fileLifecycleEvents = unix.IN_CREATE | unix.IN_DELETE
	fileWriteEvents = unix.IN_MODIFY | unix.IN_CLOSE_WRITE
	fileMoveEvents = unix.IN_MOVED_FROM | unix.IN_MOVED_TO
)

func defaultMask() uint32 {
	return fileLifecycleEvents | fileWriteEvents | fileMoveEvents
}

func (watcher *Watcher) Add(path string) error {
	if _, exists := watcher.paths[path]; exists {
		return nil
	}

	mask := defaultMask()

	watchDescriptor, err := watcher.addWatch(watcher.fileDescriptor, path, mask)

	if err != nil {
		return  err
	}

	watcher.registerDescriptor(watchDescriptor, path)
	watcher.registerPath(watchDescriptor, path)

	return nil
}

func (watcher *Watcher) Events() <-chan Event {
	return watcher.events
}

func (watcher *Watcher) Start() {
	go watcher.readLoop()
}

func (watcher *Watcher) Close () error {
	return unix.Close(watcher.fileDescriptor)
}
