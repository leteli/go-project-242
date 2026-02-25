package files

import "fmt"

var formats = []struct {
	name  string
	value float64
}{
	{name: "EB", value: 1e18},
	{name: "PB", value: 1e15},
	{name: "TB", value: 1e12},
	{name: "GB", value: 1e9},
	{name: "MB", value: 1e6},
	{name: "KB", value: 1e3},
}

func FormatSize(size int64, humanReadable bool) (string, error) {
	if size < 0 {
		return "", fmt.Errorf("size cannot be negative")
	}
	if !humanReadable || size < 1000 {
		return fmt.Sprintf("%dB", size), nil
	}
	floatSize := float64(size)

	for _, entry := range formats {
		if floatSize >= entry.value {
			if entry.value == 0 {
				return "", fmt.Errorf("cannot divide by 0")
			}
			res := floatSize / entry.value
			return fmt.Sprintf("%.1f%s", res, entry.name), nil
		}
	}
	return fmt.Sprintf("%dB", size), nil
}
