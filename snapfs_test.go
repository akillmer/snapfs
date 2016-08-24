package snapfs

import "testing"
import "fmt"

func Test_Watcher(t *testing.T) {
	if err := Watch("."); err != nil {
		t.Error(err)
		return
	}

	Ignore("(.git)")
	Update()
	Unwatch(".")
	fmt.Println(Update())
}
