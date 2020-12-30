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
	// Filter regexp patterns for file names
	Filter []string
}

// Version of a file
type Version struct {
	// Modified timestamp of file
	Modified int64
	// Hash of file version
	Hash string
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
	Modified int64
	// Old file versions
	Versions []Version
}

// BackupName of this file
func (f *File) BackupName() string {
	return f.Hash + ".tar.gz"
}

// RestorePath of this file
func (f *File) RestorePath() string {
	return f.Path + ".restore"
}
