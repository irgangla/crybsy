package crybsy

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

func createCrybsyDir(root *Root) (string, error) {
	dir, err := filepath.Abs(filepath.Join(root.Path, ".crybsy"))
	if err != nil {
		log.Println("find crybsy root failed", err)
		return "", err
	}

	err = os.MkdirAll(dir, 0775)
	if err != nil {
		log.Println("create crybsy root failed", err)
		return "", err
	}

	return dir, nil
}

// SaveRoot saves the root data
func SaveRoot(root *Root) error {
	data, err := json.Marshal(root)
	if err != nil {
		log.Println("marshal root failed", err)
		return err
	}

	dir, err := createCrybsyDir(root)
	if err != nil {
		log.Println("crybsy folder failed", err)
		return err
	}

	err = ioutil.WriteFile(filepath.Join(dir, "root.json"), data, 0775)
	if err != nil {
		log.Println("crybsy root save failed", err)
		return err
	}

	return nil
}

// LoadRoot object for given path
func LoadRoot(path string) (*Root, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		log.Println("invalid path", err)
		return nil, err
	}

	dir := filepath.Join(absPath, ".crybsy")
	data, err := ioutil.ReadFile(filepath.Join(dir, "root.json"))
	if err != nil {
		log.Println("crybsy load root failed", err)
		return nil, err
	}

	root := new(Root)
	err = json.Unmarshal(data, root)
	if err != nil {
		log.Println("invalid json data", err)
		return nil, err
	}

	return root, nil
}

// SaveFiles saves the file information
func SaveFiles(files []File, root *Root) error {
	return saveFilesAs("files.json", files, root)
}

func saveFilesAt(path string, files []File) error {
	data, err := json.Marshal(files)
	if err != nil {
		log.Println("marshal files failed", err)
		return err
	}

	err = ioutil.WriteFile(path, data, 0775)
	if err != nil {
		log.Println("crybsy files save failed", err)
		return err
	}

	return nil
}

func saveFilesAs(name string, files []File, root *Root) error {
	dir, err := createCrybsyDir(root)
	if err != nil {
		log.Println("crybsy folder failed", err)
		return err
	}

	return saveFilesAt(filepath.Join(dir, name), files)
}

// LoadFiles saves the file information
func LoadFiles(root *Root) ([]File, error) {
	return loadFilesFrom("files.json", root)
}

func loadFilesFrom(name string, root *Root) ([]File, error) {
	dir, err := createCrybsyDir(root)
	if err != nil {
		log.Println("crybsy folder failed", err)
		return nil, err
	}

	path := filepath.Join(dir, name)
	return loadFilesAt(path)
}

func loadFilesAt(path string) ([]File, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Println("crybsy load files failed", err)
		return nil, err
	}

	var files []File
	err = json.Unmarshal(data, &files)
	if err != nil {
		log.Println("invalid json data", err)
		return nil, err
	}

	return files, nil
}
