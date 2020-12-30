package crybsy

import (
	"archive/tar"
	"compress/gzip"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sync"
)

func findBackupFolder(target string, root *Root) (string, error) {
	var dir string
	var err error
	if target == "" {
		dir, err = createCrybsyDir(root)
		if err != nil {
			return "", err
		}
		dir, err = filepath.Abs(dir)
		if err != nil {
			return "", err
		}
		dir = filepath.Join(dir, "backup")
	} else {
		dir, err = filepath.Abs(target)
		if err != nil {
			return "", err
		}
	}
	err = os.MkdirAll(dir, 0775)
	if err != nil {
		return "", err
	}
	return dir, nil
}

//BackupFile create a backup of the given file
func BackupFile(file File, root *Root, target string) (string, error) {
	dir, err := findBackupFolder(target, root)
	if err != nil {
		return "", err
	}
	filename := filepath.Join(dir, file.BackupName())
	out, err := os.Create(filename)
	if err != nil {
		return "", err
	}
	defer out.Close()

	err = createArchive(file, root, out)
	if err != nil {
		return "", err
	}

	return filename, nil
}

func createArchive(file File, root *Root, buf io.Writer) error {
	gw, err := gzip.NewWriterLevel(buf, gzip.BestCompression)
	if err != nil {
		return err
	}
	defer gw.Close()
	tw := tar.NewWriter(gw)
	defer tw.Close()

	sourcepath := filepath.Join(root.Path, file.Path)
	f, err := os.Open(sourcepath)
	if err != nil {
		return err
	}
	defer f.Close()

	info, err := f.Stat()
	if err != nil {
		return err
	}

	header, err := tar.FileInfoHeader(info, info.Name())
	if err != nil {
		return err
	}
	header.Name = file.Hash
	err = tw.WriteHeader(header)
	if err != nil {
		return err
	}
	_, err = io.Copy(tw, f)
	if err != nil {
		return err
	}
	return nil
}

//RestoreFile restores a file from backup
func RestoreFile(file File, root *Root, backupPath string) error {
	backup, err := filepath.Abs(backupPath)
	if err != nil {
		return err
	}
	in, err := os.Open(backup)
	if err != nil {
		return err
	}

	targetpath := filepath.Join(root.Path, file.RestorePath())
	out, err := os.Create(targetpath)
	if err != nil {
		return err
	}
	defer out.Close()

	err = readArchive(in, out, file.Hash)
	if err != nil {
		return err
	}

	return checkAndReplace(file, root)
}

func checkAndReplace(file File, root *Root) error {
	restorepath := filepath.Join(root.Path, file.RestorePath())
	rehash, err := Hash(restorepath)
	if err != nil {
		return err
	}

	if file.Hash != rehash {
		return errors.New("restore checksum error")
	}

	targetpath := filepath.Join(root.Path, file.Path)
	err = os.Remove(targetpath)
	if err != nil {
		return err
	}

	err = os.Rename(targetpath+".restore", targetpath)
	if err != nil {
		return err
	}

	return nil
}

func readArchive(in io.Reader, out io.Writer, name string) error {
	gr, err := gzip.NewReader(in)
	if err != nil {
		return err
	}
	defer gr.Close()
	tr := tar.NewReader(gr)

	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		if hdr.Name == name {
			if _, err := io.Copy(out, tr); err != nil {
				return err
			}
			break
		}
	}
	return nil
}

// UpdateBackup backups all changed files
func UpdateBackup(root *Root, target string) ([]File, error) {
	currentFiles, err := Update(root)
	if err != nil {
		return nil, err
	}

	targetDir, err := findBackupFolder(target, root)
	if err != nil {
		return nil, err
	}

	updateQueue := make(chan File, 10)
	updatedQueue := make(chan File, 10)
	updateResult := make(chan []File, 1)

	var checkGroup sync.WaitGroup
	var updateGroup sync.WaitGroup

	checkGroup.Add(1)
	go checkFiles(currentFiles, targetDir, updateQueue, &checkGroup)

	for i := 0; i < 8; i++ {
		updateGroup.Add(1)
		go updateBackupFile(updateQueue, targetDir, root, updatedQueue, &updateGroup)
	}

	go collectUpdatedFiles(updatedQueue, updateResult)

	checkGroup.Wait()
	close(updateQueue)
	updateGroup.Wait()
	close(updatedQueue)

	return <-updateResult, nil
}

func collectUpdatedFiles(updatedQueue chan File, result chan []File) {
	files := make([]File, 0)
	for f := range updatedQueue {
		files = append(files, f)
	}
	result <- files
}

func updateBackupFile(updateQueue chan File, targetDir string, root *Root, updatedQueue chan File, wg *sync.WaitGroup) {
	defer wg.Done()

	for f := range updateQueue {
		backup, err := BackupFile(f, root, targetDir)
		if err != nil {
			log.Println("backup error", f.Path, err)
		} else {
			log.Println("backup", f.Path, "->", backup)
			updatedQueue <- f
		}
	}
}

func checkFiles(files []File, targetDir string, updateQueue chan File, wg *sync.WaitGroup) {
	defer wg.Done()

	for _, f := range files {
		if doBackup(f, targetDir) {
			updateQueue <- f
		}
	}
}

func doBackup(file File, backupFolder string) bool {
	backup := filepath.Join(backupFolder, file.BackupName())
	backupPath, err := filepath.Abs(backup)
	if err != nil {
		return true
	}

	_, err = os.Stat(backupPath)
	if err != nil {
		return true
	}

	return false
}

// VerifyBackup backups all changed files
func VerifyBackup(root *Root, target string) ([]File, []string, error) {
	currentFiles, err := Update(root)
	if err != nil {
		return nil, nil, err
	}

	targetDir, err := findBackupFolder(target, root)
	if err != nil {
		return nil, nil, err
	}

	missingBackups := make([]File, 0)
	hashMap := ByHash(currentFiles)
	for _, fs := range hashMap {
		backupPath := filepath.Join(targetDir, fs[0].BackupName())
		_, err := os.Stat(backupPath)
		if err != nil {
			missingBackups = append(missingBackups, fs...)
		}
	}

	brokenBackups := make([]string, 0)
	backups, err := ioutil.ReadDir(targetDir)
	if err != nil {
		return nil, nil, err
	}
	for _, b := range backups {
		if !b.IsDir() {
			path := filepath.Join(targetDir, b.Name())
			hash, err := nameToHash(b.Name())
			if err != nil {
				log.Println("invalid backup", path)
				continue
			}
			if !isBackupValid(path, hash) {
				brokenBackups = append(brokenBackups, hash)
			}
		}
	}

	return missingBackups, brokenBackups, nil
}

func nameToHash(name string) (string, error) {
	ending := ".tar.gz"
	le := len(ending)
	ln := len(name)
	l := ln - le

	if ln < le || name[l+1:] != ending {
		return "", errors.New("invalid file name")
	}

	bytes, err := hex.DecodeString(name[:l])
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", bytes), nil

}

func isBackupValid(path string, hash string) bool {
	in, err := os.Open(path)
	if err != nil {
		return false
	}
	out, err := ioutil.TempFile("", "crybsy_verify")
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()
	err = readArchive(in, out, hash)
	if err != nil {
		return false
	}
	h, err := Hash(filepath.Join(os.TempDir(), out.Name()))
	if err != nil {
		return false
	}
	if h == hash {
		return true
	}
	return false
}
