package watcher

import (
	"errors"
	"log"
	"path/filepath"
	"unsafe"

	"golang.org/x/sys/unix"
)


func (watcher *Watcher) readEvent(buffer []byte, offset  int) (*Event, int) {
	raw := (*unix.InotifyEvent)(unsafe.Pointer(&buffer[offset]),)
	nameLen := int(raw.Len)
	nameBytes := buffer[offset+unix.SizeofInotifyEvent : offset+unix.SizeofInotifyEvent+nameLen]

	name := string(nameBytes)

	basePath  := watcher.watches[int(raw.Wd)]

	fullPath := filepath.Join(basePath, name)

	event := watcher.parseEvent(raw, fullPath)
	nextOffset := offset + unix.SizeofInotifyEvent + nameLen
	return event, nextOffset
}

func (watcher *Watcher) readLoop() {
	buffer := make([]byte,4096)

	for {
		bytesRead, err := unix.Read(
			watcher.fileDescriptor,
			buffer,
		)

		if err != nil {
			if errors.Is(err, unix.EINTR) {
				continue
			}

			if errors.Is(err, unix.EBADF) {
				return
			}

			log.Printf("inotify read error: %v", err)

			return
		}

		offset := 0

		for offset < bytesRead {
			event, nextOffset := watcher.readEvent(buffer, offset)

			if event != nil {
				watcher.events <- *event
			}

			offset = nextOffset
		}
	}
}
