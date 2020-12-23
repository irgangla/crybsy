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

	filter := make([]string, 2)
	filter[0] = "[.]git.*"
	filter[1] = "[.]DS.*"
	root.Filter = filter

	fmt.Printf("Root: %v\n\n", root)

	files := crybsy.Collect(crybsy.Scan(root))

	fmt.Println("Files:")
	for _, f := range files {
		fmt.Println(f.Hash, f.Path)
	}

	dup := crybsy.Duplicates(crybsy.ByHash(files))
	fmt.Println("\nDuplicates:")
	for hash, files := range dup {
		fmt.Printf("%v:\n", hash)
		for i, f := range files {
			fmt.Println(i, f.Path)
		}
	}
}
