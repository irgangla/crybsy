package crybsy

import (
	"bufio"
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"os"
)

// Hash calculates a hash for the given file
func Hash(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	reader := bufio.NewReader(file)

	buffer := make([]byte, 4*1024)
	hashFunc := sha256.New()
	for {
		n, err := reader.Read(buffer)
		hashFunc.Write(buffer[:n])
		if err == io.EOF {
			break
		} else if err != nil {
			return "", err
		}
	}

	hashVal := fmt.Sprintf("%x", hashFunc.Sum(nil))
	log.Println(hashVal, path)
	return hashVal, nil
}
