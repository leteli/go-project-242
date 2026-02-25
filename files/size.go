package files

import (
	"fmt"
	"os"
)

func GetSize(path string) (int, error) {
	if path == "" {
		return 0, fmt.Errorf("file path %s is empty", path)
	}
	fileInfo, err := os.Lstat(path)
	if err != nil {
		return 0, err
	}
	if !fileInfo.IsDir() {
		return int(fileInfo.Size()), nil
	}
	// NB: 1 level of recursion only
	files, err := os.ReadDir(path)
	if err != nil {
		return 0, err
	}
	var sizeBytes int64
	for _, file := range files {
		info, err := file.Info()
		if err != nil {
			return 0, err
		}
		if !info.IsDir() {
			sizeBytes += info.Size()
		}
	}
	return int(sizeBytes), nil
}
