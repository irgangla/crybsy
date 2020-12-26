package crybsy

import (
	"log"
)

// Init folder as CryBSy root
func Init(path string) (*Root, error) {
	root, err := NewRoot(path)
	if err != nil {
		return nil, err
	}

	err = SaveRoot(root)
	if err != nil {
		return nil, err
	}

	return root, nil
}

// Load CryBSy root for given path
func Load(path string) (*Root, error) {
	root, err := LoadRoot(path)
	if err != nil {
		log.Println("path is not CryBSy root", err)
		root, err = Init(path)
		if err != nil {
			return nil, err
		}
	}
	return root, nil
}

// Update file data for root
func Update(root *Root) ([]File, error) {
	files, err := LoadFiles(root)
	if err != nil {
		log.Println("no CryBSy file data found", err)
		return Collect(Scan(root)), nil
	}

	currentFiles := Collect(Scan(root))
	return UpdateFiles(files, currentFiles), nil
}
