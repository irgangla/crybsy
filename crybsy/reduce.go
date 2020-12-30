package crybsy

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

// Reduce versions of all files to the given max
func Reduce(root *Root, max int) error {
	files, err := LoadFiles(root)
	if err != nil {
		return err
	}

	for _, f := range files {
		if f.Versions != nil && len(f.Versions) > max {
			l := len(f.Versions) - max
			f.Versions = f.Versions[l:]
		}
	}

	SaveFiles(files, root)

	return nil
}

// ReduceBackup deletes all no longer needed backup dumps
func ReduceBackup(root *Root, target string) error {
	targetDir, err := findBackupFolder(target, root)
	if err != nil {
		return err
	}

	backups, err := getBackups(root)
	if err != nil {
		return err
	}

	infos, err := ioutil.ReadDir(targetDir)
	if err != nil {
		return err
	}

	for _, info := range infos {
		hash, err := nameToHash(info.Name())
		if err != nil {
			log.Println(info.Name(), err)
			continue
		}
		_, ok := backups[hash]
		if !ok {
			path := filepath.Join(targetDir, info.Name())
			err = os.Remove(path)
			if err != nil {
				log.Println(info.Name(), err)
			}
		}
	}

	return nil
}

func getBackups(root *Root) (map[string][]File, error) {
	files, err := LoadFiles(root)
	if err != nil {
		return nil, err
	}
	backups := make(map[string][]File)
	for _, f := range files {
		addFile(f.Hash, f, &backups)
		if f.Versions != nil {
			for _, v := range f.Versions {
				addFile(v.Hash, f, &backups)
			}
		}
	}
	return backups, nil
}

func addFile(hash string, file File, fileMap *map[string][]File) {
	fm := (*fileMap)
	fs, ok := fm[hash]
	if !ok {
		fs = make([]File, 0)
	}
	fs = append(fs, file)
	fm[hash] = fs
}
