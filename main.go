package main

import (
	"fmt"
	"log"

	"github.com/crybsy/crybsy/crybsy"
)

func main() {
	fmt.Println("CryBSy")

	root, err := crybsy.NewRoot(".")
	if err != nil {
		panic(err)
	}

	log.Println(root)

	files := crybsy.Collect(crybsy.Scan(root))

	for _, f := range files {
		fmt.Println(f.Path, f.Name, f.Hash)
	}
}
