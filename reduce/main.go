package main

import (
	"fmt"
	"os"

	"github.com/crybsy/crybsy/crybsy"
)

func main() {
	fmt.Println("CryBSy Reduce\n-----")

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

	err = crybsy.Reduce(root, 3)
	if err != nil {
		panic(err)
	}

	err = crybsy.ReduceBackup(root, target)
	if err != nil {
		panic(err)
	}
}
