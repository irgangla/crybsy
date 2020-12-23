package crybsy

import (
	"os"
	"path/filepath"
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
	absDir, e := filepath.Abs(dir)
	if e != nil {
		t.Error("abs from temp dir", e)
	}
	r, e := NewRoot(dir)
	if e != nil {
		t.Error("new root from temp dir", e)
	}
	if r == nil {
		t.Fatal("new root is nil")
	}

	if r.Path != absDir {
		t.Error("root path is wrong", r.Path, dir)
	}
}

func Test_Scan(t *testing.T) {
	root, err := NewRoot("test")
	if err != nil {
		t.Fatal(err)
	}

	files, errors, wg := Scan(root)

	wg.Wait()
	close(files)
	close(errors)

	for e := range errors {
		t.Error(e)
	}

	fs := make([]File, 0)
	for f := range files {
		fs = append(fs, f)
	}

	if len(fs) != 3 {
		t.Error("wrong file length")
	}
}
