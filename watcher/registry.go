package watcher

func (watcher *Watcher) registerDescriptor(
	watchDescriptor int,
	path string,
) {
	watcher.watches[watchDescriptor] = path
}

func (watcher *Watcher) registerPath(
	watchDescriptor int,
	path string,
) {
	watcher.paths[path] = watchDescriptor
}
