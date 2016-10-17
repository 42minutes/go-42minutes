package minutes

import (
	"os"
	"path/filepath"
	"regexp"

	fsnotify "github.com/fsnotify/fsnotify"
	recwatch "github.com/xyproto/recwatch"
)

var (
	fileWatcherIgnoreRegexs = []*regexp.Regexp{
		regexp.MustCompile(`.DS_Store$`),
	}
)

// FileWatcher watches a directory and notifies Notifiees for file changes
type FileWatcher struct {
	notifiees []WatcherNotifiee
}

// Watch starts watching a folder
func (fw *FileWatcher) Watch(dir string) error {
	if err := fw.walk(dir); err != nil {
		return err
	}

	rw, err := recwatch.NewRecursiveWatcher(dir)
	if err != nil {
		return err
	}

	defer rw.Close()

	done := make(chan bool)
	go func() {
		for {
			select {
			case event := <-rw.Events:
				fp := event.Name
				if fw.blacklisted(fp) {
					continue
				}
				if event.Op&fsnotify.Create == fsnotify.Create {
					fw.notify(NotificationFileCreated, fp)
				} else if event.Op&fsnotify.Write == fsnotify.Write {
					fw.notify(NotificationFileUpdated, fp)
				} else if event.Op&fsnotify.Remove == fsnotify.Remove {
					fw.notify(NotificationFileRemoved, fp)
				}
			case err := <-rw.Errors:
				// TODO(geoah) Better error handling
				log.Info("Watcher returned an error:", err)
			}
		}
	}()

	err = rw.AddFolder(dir)
	if err != nil {
		return err
	}

	<-done

	return nil
}

// notify all notifiees about os file events
func (fw *FileWatcher) notify(ntt NotificationType, fp string) {
	for _, ntf := range fw.notifiees {
		ntf.HandleWatcherNotification(ntt, fp)
	}
}

// blacklisted files and directories
func (fw *FileWatcher) blacklisted(fp string) bool {
	if info, err := os.Stat(fp); err == nil && info.IsDir() {
		return true
	}
	for _, rx := range fileWatcherIgnoreRegexs {
		if mt := rx.MatchString(fp); mt {
			return true
		}
	}
	return false
}

// walk through a directory recursively and notifies all notifiees
// about all found files and directories as NotificationFileCreated
func (fw *FileWatcher) walk(dir string) error {
	return filepath.Walk(dir, func(fp string, f os.FileInfo, err error) error {
		if fw.blacklisted(fp) == false {
			fw.notify(NotificationFileCreated, fp)
		}
		return nil
	})
}

// Notify registers a Notifiee to be notified about file updates
func (fw *FileWatcher) Notify(wn WatcherNotifiee) {
	fw.notifiees = append(fw.notifiees, wn)
}
