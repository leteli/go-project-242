package files

import (
	"fmt"
	"os"
	"strings"
)

func isHidden(filename string) bool {
	return strings.HasPrefix(filename, ".")
}

func GetSize(path string, withHidden bool) (int64, error) {
	if path == "" {
		return 0, fmt.Errorf("file path is empty")
	}
	fileInfo, err := os.Lstat(path)
	if err != nil {
		return 0, err
	}
	hidden := isHidden(fileInfo.Name())
	if hidden && !withHidden {
		return 0, nil
	}
	if !fileInfo.IsDir() {
		return fileInfo.Size(), nil
	}
	// NB: 1 level of recursion only
	files, err := os.ReadDir(path)
	if err != nil {
		return 0, err
	}
	var sizeBytes int64
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		if isHidden(file.Name()) && !withHidden {
			continue
		}
		info, err := file.Info()
		if err != nil {
			return 0, err
		}
		sizeBytes += info.Size()
	}
	return sizeBytes, nil
}
