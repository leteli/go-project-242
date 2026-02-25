package files

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func createTestFile(testFilePath string, size int) error {
	bytes := strings.Repeat("1", size)
	err := os.WriteFile(testFilePath, []byte(bytes), 0644)
	return err
}

func TestGetPathSize(t *testing.T) {
	tempDir := t.TempDir()
	file5B := filepath.Join(tempDir, "file5.txt")
	err := createTestFile(file5B, 5)
	if err != nil {
		t.Fatalf("Unable to create test file, %v", err)
	}
	file10B := filepath.Join(tempDir, "file10.txt")
	err = createTestFile(file10B, 10)
	if err != nil {
		t.Fatalf("Unable to create test file, %v", err)
	}
	testCases := []struct {
		name         string
		path         string
		expectedSize int
		expectError  bool
	}{
		{
			name:         "1 file of 5 bytes",
			path:         file5B,
			expectedSize: 5,
			expectError:  false,
		},
		{
			name:         "1 file of 10 bytes",
			path:         file10B,
			expectedSize: 10,
			expectError:  false,
		},
		{
			name:         "1 directory with 2 files",
			path:         tempDir,
			expectedSize: 15,
			expectError:  false,
		},
		{
			name:         "1 directory with 2 files",
			path:         "",
			expectedSize: 0,
			expectError:  true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			size, err := GetSize(testCase.path)
			hasError := err != nil
			if hasError != testCase.expectError {
				t.Fatalf("Error expected: %v; received: %v", testCase.expectError, err)
			}
			if size != testCase.expectedSize {
				t.Errorf("Expected size: %d; actual size %d", testCase.expectedSize, size)
			}
		})
	}
}
