package minutes

import (
	"log"

	fsnotify "github.com/fsnotify/fsnotify"
	recwatch "github.com/xyproto/recwatch"
)

// FileWatcher watches a directory and notifies Notifiees for file changes
type FileWatcher struct {
	notifiees []WatcherNotifiee
}

// Watch starts watching a folder
func (fw *FileWatcher) Watch(path string, recursive bool) error {
	rw, err := recwatch.NewRecursiveWatcher(path)
	if err != nil {
		return err
	}

	defer rw.Close()

	done := make(chan bool)
	go func() {
		for {
			select {
			case event := <-rw.Events:
				log.Println("event:", event)
				if event.Op&fsnotify.Create == fsnotify.Create {
					fw.notify(NotificationFileCreated, event.Name)
				} else if event.Op&fsnotify.Write == fsnotify.Write {
					fw.notify(NotificationFileUpdated, event.Name)
				} else if event.Op&fsnotify.Remove == fsnotify.Remove {
					fw.notify(NotificationFileRemoved, event.Name)
				}
			case err := <-rw.Errors:
				log.Println("Watcher returned an error:", err) // TODO(geoah) Better error handling
			}
		}
	}()

	err = rw.AddFolder(path)
	if err != nil {
		return err
	}

	<-done

	return nil
}

func (fw *FileWatcher) notify(ntt NotificationType, path string) {
	for _, ntf := range fw.notifiees {
		ntf.HandleWatcherNotification(ntt, path)
	}
}

// Notify registers a Notifiee to be notified about file updates
func (fw *FileWatcher) Notify(wn WatcherNotifiee) {
	fw.notifiees = append(fw.notifiees, wn)
}
