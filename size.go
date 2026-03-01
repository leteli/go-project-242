package code

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
)

var ErrorEmptyPath = errors.New("file path is empty")

func isHidden(filename string) bool {
	return strings.HasPrefix(filename, ".")
}

func getDirSize(path string, withHidden bool, recursive bool) (int64, error) {
	files, err := os.ReadDir(path)
	if err != nil {
		return 0, err
	}
	var sizeBytes int64
	for _, file := range files {
		if isHidden(file.Name()) && !withHidden {
			continue
		}
		if file.IsDir() {
			if !recursive {
				continue
			}
			dirPath := filepath.Join(path, file.Name())
			dirSize, err := getDirSize(dirPath, withHidden, recursive)
			if err != nil {
				return 0, err
			}
			sizeBytes += dirSize
		} else {
			info, err := file.Info()
			if err != nil {
				return 0, err
			}
			sizeBytes += info.Size()
		}
	}
	return sizeBytes, nil
}

func GetPathSize(path string, withHidden bool, recursive bool) (int64, error) {
	if path == "" {
		return 0, ErrorEmptyPath
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
	sizeBytes, err := getDirSize(path, withHidden, recursive)
	if err != nil {
		return 0, err
	}
	return sizeBytes, nil
}
