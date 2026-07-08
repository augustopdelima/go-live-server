package watcher

import (
	"io/fs"
	"path/filepath"
)


func WalkAndWatch(watcher *Watcher, root string, matcher *IgnoreMatcher) error {
	return filepath.WalkDir(root, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !entry.IsDir() {
			return nil
		}

		if path != root && matcher.Match(path) {

			return filepath.SkipDir
		}

		return watcher.Add(path)

	})
}
