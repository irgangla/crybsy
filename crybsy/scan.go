package crybsy

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"regexp"
	"sync"
)

type scanner struct {
	Root      *Root
	Files     chan File
	Errors    chan error
	WaitGroup *sync.WaitGroup
	Filter    []*regexp.Regexp
}

// NewRoot creates a new CryBSy Root
func NewRoot(path string) (*Root, error) {
	if len(path) == 0 {
		return nil, errors.New("empty path is not valid")
	}

	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	if !info.IsDir() {
		return nil, errors.New("path is not a directory")
	}

	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}
	root := new(Root)
	root.Path = absPath
	root.Host, err = os.Hostname()
	if err != nil {
		root.Host = "unknown"
	}
	user, err := user.Current()
	if err == nil {
		root.User.Name = user.Name
		root.User.UID = user.Uid
		root.User.GID = user.Gid
	}
	root.ID = calculateRootID(root)

	return root, nil
}

func calculateRootID(root *Root) string {
	hashFunc := sha256.New()
	hashFunc.Write([]byte(root.Path))
	hashFunc.Write([]byte(root.Host))
	hashFunc.Write([]byte(root.User.Name))
	hashFunc.Write([]byte(root.User.GID))
	hashFunc.Write([]byte(root.User.UID))
	return fmt.Sprintf("%x", hashFunc.Sum(nil))
}

// Collect all files found by a scan
func Collect(files chan File, errors chan error, wg *sync.WaitGroup) []File {
	wg.Wait()
	close(errors)
	close(files)

	for e := range errors {
		log.Println("CryBSy: collect files", e)
	}

	fs := make([]File, 0)
	for f := range files {
		fs = append(fs, f)
	}
	return fs
}

// Scan the root tree for files
func Scan(root *Root) (chan File, chan error, *sync.WaitGroup) {
	var wg sync.WaitGroup
	errors := make(chan error, 1000)

	patterns := make([]*regexp.Regexp, 0)
	if root.Filter != nil {
		for _, f := range root.Filter {
			regexp, err := regexp.Compile(f)
			if err != nil {
				errors <- err
			} else {
				patterns = append(patterns, regexp)
			}
		}
	}

	scan := scanner{
		Root:      root,
		Files:     make(chan File, 1000),
		Errors:    errors,
		WaitGroup: &wg,
		Filter:    patterns,
	}

	wg.Add(1)
	go scanRecursive(root.Path, scan)

	return scan.Files, scan.Errors, scan.WaitGroup
}

func scanRecursive(path string, scan scanner) {
	defer scan.WaitGroup.Done()
	callback := func(filePath string, file os.FileInfo, err error) error {
		if err != nil {
			scan.Errors <- err
			return err
		}

		if !file.IsDir() {
			scan.WaitGroup.Add(1)
			go handleFile(filePath, file, scan)
		}

		return err
	}
	err := filepath.Walk(path, callback)
	if err != nil {
		scan.Errors <- err
	}
}

func filterFile(path string, scan scanner) bool {
	for _, exp := range scan.Filter {
		if exp.Match([]byte(path)) {
			return true
		}
	}
	return false
}

func handleFile(path string, file os.FileInfo, scan scanner) {
	defer scan.WaitGroup.Done()

	if filterFile(path, scan) {
		return
	}

	absPath, err := filepath.Abs(path)
	if err != nil {
		scan.Errors <- err
		return
	}
	hash, err := Hash(absPath)
	if err != nil {
		scan.Errors <- err
		return
	}

	relPath, err := filepath.Rel(scan.Root.Path, absPath)
	if err != nil {
		scan.Errors <- err
		return
	}
	modified := file.ModTime().Unix()
	_, name := filepath.Split(path)
	f := File{
		Path:     relPath,
		Name:     name,
		RootID:   scan.Root.ID,
		Modified: modified,
		Hash:     hash,
		FileID:   hash,
	}
	scan.Files <- f
}
