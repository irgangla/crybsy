package crybsy

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func Test_Backup(t *testing.T) {
	root, err := NewRoot("test")
	if err != nil {
		t.Fatal(err)
	}

	path := "subfolder/s1.txt"
	hash, err := Hash(filepath.Join(root.Path, path))
	if err != nil {
		t.Fatal(err)
	}
	var file = File{
		Path: path,
		Name: "s1.txt",
		Hash: hash,
	}

	backupPath, err := BackupFile(file, root, "")
	if err != nil {
		t.Fatal(err)
	}

	_, err = os.Stat(backupPath)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_Restore(t *testing.T) {
	root, err := NewRoot("test")
	if err != nil {
		t.Fatal(err)
	}

	path := "subfolder/s1.txt"
	hash, err := Hash(filepath.Join(root.Path, path))
	if err != nil {
		t.Fatal(err)
	}
	var file = File{
		Path: path,
		Name: "s1.txt",
		Hash: hash,
	}

	backupPath, err := BackupFile(file, root, "")
	if err != nil {
		t.Fatal(err)
	}

	err = RestoreFile(file, root, backupPath)
	if err != nil {
		t.Fatal(err)
	}

	data, err := ioutil.ReadFile(filepath.Join(root.Path, path))
	if err != nil {
		t.Fatal(err)
	}

	if string(data) != "Hallo, Welt!\n" {
		t.Error("wrong file content", string(data))
	}
}
