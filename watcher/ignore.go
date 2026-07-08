package watcher

import (
	"path/filepath"
	"strings"
)

const hidden_prefix = "."

type IgnoreMatcher struct {
	directories map[string]struct{}
}

func NewIgnoreMatcher() *IgnoreMatcher {
	return &IgnoreMatcher{
		directories: map[string]struct{}{
			"node_modules":{},
			"vendor":{},
			".git":{},
			"dist":{},
		},
	}
}

func (m *IgnoreMatcher) Match(path string) bool {
	name := filepath.Base(path)

	if strings.HasPrefix(name,".") {
		return true
	}

	_, ignored := m.directories[name]

	return ignored
}
