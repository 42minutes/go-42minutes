package minutes

// NotificationType for the types of Watcher notifications
type NotificationType int

const (
	// NotificationFileCreated - file created
	NotificationFileCreated NotificationType = iota
	// NotificationFileUpdated - file updated (optional)
	NotificationFileUpdated
	// NotificationFileMoved - file moved (optiona)
	NotificationFileMoved
	// NotificationFileRemoved - file removed
	NotificationFileRemoved
)

// Watcher watches a directory or api notifies Notifiees for file changes
type Watcher interface {
	// Watch starts watching a directory or api
	// or errors with ErrInternalServer
	Watch(path string, recursive bool) error
	// Notify adds Notifiees that want to be notified about file changes
	Notify(WatcherNotifiee)
}

// WatcherNotifiee is anything that wants to be notified about file changes
type WatcherNotifiee interface {
	// HandleWatcherNotification handles Watcher notifications
	HandleWatcherNotification(notifType NotificationType, path string)
}
