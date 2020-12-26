package crybsy

import (
	"fmt"
)

// PrintRoot prints all root information fromated to stdout
func PrintRoot(root *Root) {
	fmt.Println("CryBSy Root:")
	fmt.Println("- Path:", root.Path)
	fmt.Println("- ID:", root.ID)
	fmt.Println("- Host:", root.Host)
	fmt.Printf("- User: %v (%v, %v)", root.User.Name, root.User.UID, root.User.GID)
	fmt.Printf("- Filter: ")
	if root.Filter == nil {
		fmt.Println("no filter")
	} else {
		for _, f := range root.Filter {
			fmt.Printf("\"%v\" ", f)
		}
	}
	fmt.Printf("\n-----\n\n")
}
