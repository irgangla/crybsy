package main

import (
	"fmt"

	"github.com/crybsy/crybsy/crybsy"
)

func main() {
	fmt.Println("CryBSy")

	root, err := crybsy.NewRoot(".")
	if err != nil {
		panic(err)
	}
	crybsy.SetDefaultFilter(root)
	crybsy.SaveRoot(root)

	crybsy.PrintRoot(root)

	files, err := crybsy.Update(root)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Files (%v):\n", len(files))
	for _, f := range files {
		fmt.Println(f.Hash, f.Path)
	}

	crybsy.SaveFiles(files, root)

	dup := crybsy.Duplicates(crybsy.ByHash(files))
	fmt.Println("\nDuplicates:")
	for hash, files := range dup {
		fmt.Printf("%v:\n", hash)
		for i, f := range files {
			fmt.Println(i, f.Path)
		}
	}
}
