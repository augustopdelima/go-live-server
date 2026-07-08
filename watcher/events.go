package watcher

import "time"

type EventType string

const (
	Create EventType = "create"
	Delete EventType = "delete"
	Update EventType = "update"
	Move EventType = "move"
)

type Event struct {
	Type EventType
	Path string
	IsDir bool
	Time time.Time
}
