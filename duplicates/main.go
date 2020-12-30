package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/crybsy/crybsy/crybsy"
)

func main() {
	fmt.Println("CryBSy Duplicates\n-----")

	args := os.Args
	path := "."
	out := "./duplicates.json"
	var err error

	if len(args) > 1 {
		path = args[1]
	}
	path, err = filepath.Abs(path)
	if err != nil {
		log.Panicln(err)
	}

	if len(args) > 2 {
		out = args[2]
	}
	out, err = filepath.Abs(out)
	if err != nil {
		log.Panicln(err)
	}

	fmt.Println("Path:", path)
	fmt.Println("Out:", out)
	fmt.Println("-----")

	root, err := crybsy.NewRoot(path)
	if err != nil {
		panic(err)
	}
	crybsy.SetDefaultFilter(root)
	crybsy.SaveRoot(root)

	crybsy.PrintRoot(root)

	files, err := crybsy.LoadFiles(root)
	if err != nil {
		panic(err)
	}

	fmt.Println("Files:", len(files))

	dup := crybsy.Duplicates(crybsy.ByHash(files))

	saveDuplicates(dup, out)
	printDuplicates(dup)
}

func printDuplicates(dup map[string][]crybsy.File) {
	fmt.Println("\nDuplicates:", len(dup))
	for hash, files := range dup {
		fmt.Printf("%v: ", hash)
		for _, f := range files {
			fmt.Printf("%v ", f.Path)
		}
	}
	fmt.Println()
}

func saveDuplicates(dup map[string][]crybsy.File, out string) {
	data, err := json.Marshal(dup)
	if err != nil {
		log.Panicln("marshal duplicates failed", err)
	}

	err = ioutil.WriteFile(out, data, 0775)
	if err != nil {
		log.Panicln("save duplicates failed", err)
	}
}
