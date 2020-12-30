package main

import (
	"fmt"
	"os"

	"github.com/crybsy/crybsy/crybsy"
)

func main() {
	fmt.Println("CryBSy Scanner\n-----")

	args := os.Args
	path := "."

	if len(args) > 1 {
		path = args[1]
	}

	fmt.Println("Path:", path)
	fmt.Println("-----")

	root, err := crybsy.NewRoot(path)
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

	fmt.Println("Files:", len(files))

	crybsy.SaveFiles(files, root)
}
