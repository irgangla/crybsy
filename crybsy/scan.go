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
	"time"
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

// SetDefaultFilter for root
func SetDefaultFilter(root *Root) {
	if root.Filter == nil {
		root.Filter = make([]string, 0)
	}
	root.Filter = append(root.Filter, "[.]git.*")
	root.Filter = append(root.Filter, "[.]DS.*")
	root.Filter = append(root.Filter, "[.]crybsy.*")
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

// Scan the root tree for files
func Scan(root *Root) (chan File, chan error, *sync.WaitGroup) {
	var wg sync.WaitGroup
	errors := make(chan error, 10000)

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
		Files:     make(chan File, 100),
		Errors:    errors,
		WaitGroup: &wg,
		Filter:    patterns,
	}

	wg.Add(1)
	go scanRecursive(root.Path, scan)

	return scan.Files, scan.Errors, scan.WaitGroup
}

func scanRecursive(path string, scan scanner) {
	log.Println("Scan folder", path)
	defer scan.WaitGroup.Done()
	callback := func(filePath string, file os.FileInfo, err error) error {
		if err != nil {
			scan.Errors <- err
			return err
		}

		if !file.IsDir() {
			handleFile(filePath, file, scan)
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
	if filterFile(path, scan) {
		return
	}

	absPath, err := filepath.Abs(path)
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
	}
	scan.Files <- f
}

// UpdateFiles merge the old and new scan
func UpdateFiles(oldFiles []File, root *Root, files chan File, errors chan error, wg *sync.WaitGroup) []File {
	log.Println("Update file list...")

	start := time.Now().UnixNano()
	fileMap := ByPath(oldFiles)
	end := time.Now().UnixNano()
	delta := end - start
	log.Println("Mop old files:", (delta / 1000000), "ms")

	start = time.Now().UnixNano()
	updatedFiles := make(chan File, 100)
	var wg2 sync.WaitGroup
	wg2.Add(1)
	go logErrors(errors, &wg2)
	for i := 0; i < 8; i++ {
		wg2.Add(1)
		go updateFile(files, root, fileMap, updatedFiles, &wg2)
	}

	fileList := make(chan []File, 1)
	go collectFiles(updatedFiles, fileList)

	wg.Wait()
	close(files)
	close(errors)
	end = time.Now().UnixNano()
	delta = end - start
	log.Println("Disk files scanned:", (delta / 1000000), "ms")

	start = time.Now().UnixNano()
	wg2.Wait()
	close(updatedFiles)
	end = time.Now().UnixNano()
	delta = end - start
	log.Println("Process files:", (delta / 1000000), "ms")

	return <-fileList
}

func collectFiles(files chan File, res chan []File) {
	fileList := make([]File, 0)
	for f := range files {
		fileList = append(fileList, f)
	}
	res <- fileList
}

func logErrors(errors chan error, wg *sync.WaitGroup) {
	defer wg.Done()
	for err := range errors {
		log.Println("scan file error", err)
	}
}

func updateFile(files chan File, root *Root, fileMap map[string]File, updateFiles chan File, wg *sync.WaitGroup) {
	defer wg.Done()
	for f := range files {
		start := time.Now().UnixNano()
		of, ok := fileMap[f.Path]
		if !ok {
			// new file found
			hash, err := Hash(f.Path)
			if err != nil {
				log.Println("file hash error", err)
			} else {
				f.FileID = hash
				f.Hash = hash
				updateFiles <- f
			}
		} else {
			// handle old file
			if of.Modified == f.Modified {
				updateFiles <- of
			} else {
				hash, err := Hash(filepath.Join(root.Path, f.Path))
				if err != nil {
					log.Println("file hash error", err)
					updateFiles <- of
				} else {
					v := Version{
						Modified: of.Modified,
						Hash:     of.Hash,
					}
					of.Versions = append(of.Versions, v)
					of.Hash = hash
					updateFiles <- of
				}
			}
		}

		end := time.Now().UnixNano()
		delta := end - start
		log.Println("Update file", f.Path, "Time:", (delta / 1000000), "ms")
	}
}
