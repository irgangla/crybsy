package main

import (
	"fmt"
	"os"

	"github.com/crybsy/crybsy/crybsy"
)

func main() {
	fmt.Println("CryBSy Verify Backup\n-----")

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

	missing, broken, err := crybsy.VerifyBackup(root, target)
	if err != nil {
		panic(err)
	}

	if len(missing) > 0 {
		fmt.Println("Missing files:")
		for _, f := range missing {
			fmt.Println(f.Path)
		}
	}

	if len(broken) > 0 {
		fmt.Println("Broken backups:")
		for _, b := range broken {
			fmt.Println(b)
		}
	}

	if len(broken) == 0 && len(missing) == 0 {
		fmt.Println("Backup successful verified.")
	}
}
