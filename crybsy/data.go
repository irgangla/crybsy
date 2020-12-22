package crybsy

// Owner of the files
type Owner struct {
	// Name of the system user
	Name string
	// Uid of the user
	UID string
	// Gid if the user
	GID string
}

// Root folder of tracked folder tree
type Root struct {
	// Path to root folder on host
	Path string
	// Host name
	Host string
	// User name on the host
	User Owner
	// ID is a global unique id for tree replication
	ID string
}

// File groups all information about a specific file
type File struct {
	// Path relative to root folder
	Path string
	// Name of the file including extension
	Name string
	// Type of the file, for metadata processing
	Type string
	// MetaData of the file, e.g. extracted EXIF data
	MetaData interface{}
	// RootID is the ID of the root folder of this file
	RootID string
	// FileID is a global unique ID of this file
	FileID string
	// Hash of last file scan
	Hash string
	// Modified date of file data
	Modified uint64
}
