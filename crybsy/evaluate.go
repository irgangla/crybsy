package crybsy

// ByHash groups the files by file hash values
func ByHash(files []File) map[string][]File {
	fileMap := make(map[string][]File)
	for _, f := range files {
		list, ok := fileMap[f.Hash]
		if !ok {
			list = make([]File, 0)
		}
		list = append(list, f)
		fileMap[f.Hash] = list
	}
	return fileMap
}

// Duplicates finds files with same hash
func Duplicates(byHash map[string][]File) map[string][]File {
	fileMap := make(map[string][]File)
	for hash, files := range byHash {
		if len(files) > 1 {
			fileMap[hash] = files
		}
	}
	return fileMap
}
