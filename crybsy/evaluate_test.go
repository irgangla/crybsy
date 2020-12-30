package crybsy

import "testing"

func Test_ByHash(t *testing.T) {
	root, err := NewRoot("test")
	if err != nil {
		t.Fatal(err)
	}

	files, err := Update(root)
	if err != nil {
		t.Fatal(err)
	}

	fileMap := ByHash(files)

	for hash, files := range fileMap {
		for i, f := range files {
			if f.Hash != hash {
				t.Error(i, "wrong hash", f.Path)
			}
		}
	}
}

func Test_Duplicate(t *testing.T) {
	root, err := NewRoot("test")
	if err != nil {
		t.Fatal(err)
	}

	files, err := Update(root)
	if err != nil {
		t.Fatal(err)
	}

	fileMap := Duplicates(ByHash(files))

	for _, files := range fileMap {
		if len(files) <= 1 {
			t.Error("no duplicates", files)
		}
	}
}
