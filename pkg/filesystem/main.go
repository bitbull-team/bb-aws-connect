package filesystemlib

import (
	"os"
)

// FileExist check if path is a regular file end exist
func FileExist(filepath string) bool {
	info, err := os.Stat(filepath)
	if os.IsNotExist(err) || info.IsDir() {
		return false
	}
	return true
}

// DirectoryExist check if path is a regular directory end exist
func DirectoryExist(filepath string) bool {
	info, err := os.Stat(filepath)
	if os.IsNotExist(err) || !info.IsDir() {
		return false
	}
	return true
}
