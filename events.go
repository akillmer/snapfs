package snapfs

import (
	"os"
	"path/filepath"
)

// The Events struct contains the paths to files or directories that have been
// added, modified, or removed since the last snapfs.Update() call.
type Events struct {
	Files  subEvents        `json:"files"`
	Dirs   subEvents        `json:"dirs"`
	Errors map[string]error `json:"errors"`
}

type subEvents struct {
	Added    []string `json:"added"`
	Modified []string `json:"modified"`
	Removed  []string `json:"removed"`
}

func newEvents() Events {
	return Events{
		Errors: make(map[string]error),
	}
}

func (e *Events) error(path string, err error) {
	e.Errors[path] = err
}

func (e *Events) append(events Events) {
	e.Files.Added = append(e.Files.Added, events.Files.Added...)
	e.Files.Modified = append(e.Files.Modified, events.Files.Modified...)
	e.Files.Removed = append(e.Files.Removed, events.Files.Removed...)
	e.Dirs.Added = append(e.Dirs.Added, events.Dirs.Added...)
	e.Dirs.Modified = append(e.Dirs.Modified, events.Dirs.Modified...)
	e.Dirs.Removed = append(e.Dirs.Removed, events.Dirs.Removed...)
}

func (e *Events) add(path string, file os.FileInfo) {
	pathToFile := filepath.Join(path, file.Name())
	if file.IsDir() {
		e.Dirs.Added = append(e.Dirs.Added, pathToFile)
	} else {
		e.Files.Added = append(e.Files.Added, pathToFile)
	}
}

func (e *Events) modify(path string, file os.FileInfo) {
	pathToFile := filepath.Join(path, file.Name())
	if file.IsDir() {
		e.Dirs.Modified = append(e.Dirs.Modified, pathToFile)
	} else {
		e.Files.Modified = append(e.Files.Modified, pathToFile)
	}
}

func (e *Events) removeDir(path string, fileName string) {
	pathToFile := filepath.Join(path, fileName)
	e.Dirs.Removed = append(e.Dirs.Removed, pathToFile)
}

func (e *Events) removeFile(path string, fileName string) {
	pathToFile := filepath.Join(path, fileName)
	e.Files.Removed = append(e.Files.Removed, pathToFile)
}
