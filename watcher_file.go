package minutes

// FileWatcher watches a directory and notifies Notifiees for file changes
type FileWatcher struct{}

// Watch starts watching a folder
func (fw *FileWatcher) Watch(path string, recursive bool) error {
	return nil
}

// Notify registers a Notifiee to be notified about file updates
func (fw *FileWatcher) Notify(WatcherNotifiee) {}
