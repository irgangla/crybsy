package crybsy

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"os"
	"os/user"
)

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

	root := new(Root)
	root.Path = path
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
