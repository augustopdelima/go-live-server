package watcher

import (
	"time"

	"golang.org/x/sys/unix"
)



func (watcher *Watcher) parseEvent(
	raw *unix.InotifyEvent,
	path string,
) *Event {
	var eventType EventType

	switch {
		case raw.Mask&unix.IN_CREATE != 0:
			eventType = Create

		case raw.Mask&unix.IN_CLOSE_WRITE != 0:
			eventType = Update

		case raw.Mask&unix.IN_DELETE != 0:
			eventType = Delete

		case raw.Mask&unix.IN_MOVED_TO != 0:
			eventType = Move

		default:
		 	return nil
	}

	return &Event {
		Type: eventType,
		Path: path,
		IsDir: raw.Mask&unix.IN_ISDIR != 0,
		Time: time.Now(),
	}
}
