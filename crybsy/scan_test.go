package crybsy

import (
	"os"
	"testing"
)

func Test_NewRoot_Empty(t *testing.T) {
	_, e := NewRoot("")
	if e == nil {
		t.Error("empty root path was accepted")
	}
}
func Test_NewRoot(t *testing.T) {
	dir := os.TempDir()
	r, e := NewRoot(dir)
	if e != nil {
		t.Error("new root from temp dir", e)
	}
	if r == nil {
		t.Fatal("new root is nil")
	}

	if r.Path != dir {
		t.Error("root path is wrong", r.Path, dir)
	}
}
