package snapfs

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
)

var fs = watcher{Paths: make(map[string]snapdir)}
var fsEvents = newEvents()

type watcher struct {
	sync.Mutex
	Paths map[string]snapdir `json:"paths"`
	regex []*regexp.Regexp
}

// Watch adds a path to be monitored for any changes to its files and subdirectories.
func Watch(path string) error {
	path = filepath.Clean(path)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return err
	}
	if _, exists := fs.Paths[path]; exists {
		return fmt.Errorf("Path %s is already being watched", path)
	}
	fs.Paths[path] = newSnapdir(path)
	return nil
}

// Unwatch removes any paths that are currently being monitored for changes,
// including its subdirectories.
func Unwatch(path string) {
	for k, v := range fs.Paths {
		if strings.HasPrefix(k, path) {
			for name := range v.Dirs {
				fsEvents.removeDir(k, name)
			}
			for name := range v.Files {
				fsEvents.removeFile(k, name)
			}
			delete(fs.Paths, k)
			fsEvents.removeDir(path, k)
		}
	}
}

// Ignore uses regular expressions to determine which files or folders to ignore,
// based on file name.
func Ignore(regex ...string) {
	for _, r := range regex {
		exp := regexp.MustCompile(r)
		fs.regex = append(fs.regex, exp)
	}
}

// Snapshot returns all the currently watched paths, files, and subdirectories
// as a []byte array that has been marshaled to JSON format.
func Snapshot() ([]byte, error) {
	snapshot, err := json.Marshal(&fs.Paths)
	if err != nil {
		return []byte{}, err
	}
	return snapshot, nil
}

// Restore takes a previous Snapshot and uses it as snapfs' current state.
func Restore(snapshot []byte) error {
	restorePaths := make(map[string]snapdir)
	err := json.Unmarshal(snapshot, &restorePaths)
	if err != nil {
		return err
	}
	fs.Paths = restorePaths
	fsEvents = newEvents()
	return nil
}

func ignorePath(path string) bool {
	// no mutex locking here, since according to the golang docs:
	// A Regexp is safe for concurrent use by multiple goroutines.
	for _, r := range fs.regex {
		if r.MatchString(path) {
			return true
		}
	}
	return false
}

// Update scans all watched paths and detects any changes to files
// and subdirectories. These events are returned as a snapfs.Events struct.
func Update() Events {
	e := Events{}

	// TODO: goroutines would be awesome here, with a worker pool.
	for _, dir := range fs.Paths {
		e.append(dir.update())
	}

	// fsEvents is modified by Unwatch(), when the user removes a path
	e.append(fsEvents)
	fsEvents = newEvents()

	return e
}

func (w *watcher) set(path string, dir *snapdir) {
	// this is the only function that needs to protect
	// itself from races -- this is in place for adding goroutines later
	w.Lock()
	defer w.Unlock()
	w.Paths[path] = *dir
}
