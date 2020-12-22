package main

import (
	"fmt"
	"log"
	"os"

	"github.com/crybsy/crybsy/crybsy"
)

func main() {
	fmt.Println("CryBSy")

	root, err := crybsy.NewRoot(os.TempDir())
	if err != nil {
		panic(err)
	}

	log.Println(root)
}
