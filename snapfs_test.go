package snapfs

import (
	"strings"
	"testing"
)

func Test_Watcher(t *testing.T) {
	if err := Watch("."); err != nil {
		t.Error(err)
	}
	if Update().Dirs.Added[0] != ".git" {
		t.Error("Subdirectory `.git` should have been added")
	}
	Unwatch(".")
	if Update().Dirs.Removed[0] != ".git" {
		t.Error("Subdirectory `.git` should have been removed")
	}
}

func Test_Ignore(t *testing.T) {
	if err := Watch("."); err != nil {
		t.Error(err)
	}
	Ignore("(.go)")
	for _, v := range Update().Files.Added {
		if strings.HasSuffix(v, ".go") {
			t.Error("All .go files should be ignored")
		}
	}
	Ignore("(.md)")
	if Update().Files.Removed[0] != "README.md" {
		t.Error("README.md should no longer be watched")
	}
}

func Test_Snapshot(t *testing.T) {
	// At this point, only `LICENSE` and subdir `.git` is being watched
	snapshot, err := Snapshot()
	if err != nil {
		t.Error(err)
	}
	Unwatch(".")
	Update() // now nothing is being watched
	if err = Restore(snapshot); err != nil {
		t.Error(err)
	}
	if _, exists := fs.Paths["."]; exists == false {
		t.Error("Unable to restore snapshot")
	}
}
