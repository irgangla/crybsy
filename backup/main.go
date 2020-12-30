package main

import (
	"fmt"
	"os"

	"github.com/crybsy/crybsy/crybsy"
)

func main() {
	fmt.Println("CryBSy Backup\n-----")

	args := os.Args
	path := "."
	target := ""

	if len(args) > 1 {
		path = args[1]
	}

	if len(args) > 2 {
		target = args[2]
	}

	fmt.Println("Path:", path)
	fmt.Println("Target:", target)
	fmt.Println("-----")

	root, err := crybsy.LoadRoot(path)
	if err != nil {
		panic(err)
	}
	crybsy.PrintRoot(root)

	files, err := crybsy.UpdateBackup(root, target)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Updated files (%v):\n", len(files))
	for _, f := range files {
		fmt.Println(f.Path)
	}
}
