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

	filter := make([]string, 1)
	filter[0] = "[.]git.*"
	root.Filter = filter

	fmt.Printf("Root: %v\n\n", root)

	files := crybsy.Collect(crybsy.Scan(root))

	for _, f := range files {
		fmt.Println(f.Hash, f.Path)
	}
}
