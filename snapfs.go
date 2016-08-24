package snapfs

import (
	"io/ioutil"
	"os"
	"time"
)

type snapdir struct {
	Path  string              `json:"path"`
	Files map[string]snapfile `json:"files"`
	Dirs  map[string]snapfile `json:"dirs"`
}

type snapfile struct {
	Size    int64     `json:"size"`
	ModTime time.Time `json:"modtime"`
}

func newSnapdir(path string) snapdir {
	return snapdir{
		Path:  path,
		Files: make(map[string]snapfile),
		Dirs:  make(map[string]snapfile),
	}
}

func newSnapfile(f os.FileInfo) snapfile {
	return snapfile{Size: f.Size(), ModTime: f.ModTime()}
}

func (s *snapdir) update() Events {
	e := newEvents()
	latest := newSnapdir(s.Path)

	files, err := ioutil.ReadDir(s.Path)
	if err != nil {
		e.error(s.Path, err)
		return e
	}

	for _, f := range files {
		if ignorePath(f.Name()) {
			continue
		}
		latest.add(f)
		if exists, modified := s.contains(f); exists && modified {
			e.modify(s.Path, f)
		} else if exists == false {
			e.add(s.Path, f)
		}
	}

	// anything present in *snapdir that's missing in latest
	// has been removed
	for k := range s.Dirs {
		if _, exists := latest.Dirs[k]; exists == false {
			e.removeDir(s.Path, k)
		}
	}

	for k := range s.Files {
		if _, exists := latest.Files[k]; exists == false {
			e.removeFile(s.Path, k)
		}
	}

	fs.set(s.Path, &latest)
	return e
}

func (s *snapdir) add(f os.FileInfo) {
	if f.IsDir() {
		s.Dirs[f.Name()] = newSnapfile(f)
	} else {
		s.Files[f.Name()] = newSnapfile(f)
	}
}

func (s *snapdir) contains(f os.FileInfo) (bool, bool) {
	var (
		exists   bool
		modified bool
		snap     snapfile
	)

	if f.IsDir() {
		snap, exists = s.Dirs[f.Name()]
	} else {
		snap, exists = s.Files[f.Name()]
	}

	if exists {
		if f.ModTime() != snap.ModTime || f.Size() != snap.Size {
			modified = true
		}
	}

	return exists, modified
}
