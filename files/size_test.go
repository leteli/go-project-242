package files

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
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
		t.Fatalf("unable to create test file, %v", err)
	}
	file10B := filepath.Join(tempDir, "file10.txt")
	err = createTestFile(file10B, 10)
	if err != nil {
		t.Fatalf("unable to create test file, %v", err)
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
			if testCase.expectError {
				require.EqualError(t, err, "file path is empty")
			} else {
				require.NoError(t, err)
			}
			require.Equal(t, testCase.expectedSize, size)
		})
	}
}

func TestGetFormatSize(t *testing.T) {
	testCases := []struct {
		name           string
		size           float64
		humanReadable  bool
		expectedResult string
		expectError    bool
	}{
		{
			name:           "Size in bytes",
			size:           23711,
			humanReadable:  false,
			expectedResult: "23711B",
			expectError:    false,
		},
		{
			name:           "Human readable size in bytes",
			size:           670,
			humanReadable:  true,
			expectedResult: "670B",
			expectError:    false,
		},
		{
			name:           "Human readable size in KBs",
			size:           11_000,
			humanReadable:  true,
			expectedResult: "11.0KB",
			expectError:    false,
		},
		{
			name:           "Human readable size in MBs",
			size:           120_780_000,
			humanReadable:  true,
			expectedResult: "120.8MB",
			expectError:    false,
		},
		{
			name:           "Human readable size in GBs",
			size:           280_350_000_000,
			humanReadable:  true,
			expectedResult: "280.4GB",
			expectError:    false,
		},
		{
			name:           "Human readable size in TBs",
			size:           641_631_200_000_000,
			humanReadable:  true,
			expectedResult: "641.6TB",
			expectError:    false,
		},
		{
			name:           "Human readable size in PBs",
			size:           129_940_000_000_000_000,
			humanReadable:  true,
			expectedResult: "129.9PB",
			expectError:    false,
		},
		{
			name:           "Human readable size in EBs",
			size:           1_260_000_000_000_000_000,
			humanReadable:  true,
			expectedResult: "1.3EB",
			expectError:    false,
		},
		{
			name:           "Negative size value",
			size:           -1000,
			humanReadable:  true,
			expectedResult: "",
			expectError:    true,
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			formattedSize, err := FormatSize(int(testCase.size), testCase.humanReadable)
			if testCase.expectError {
				require.EqualError(t, err, "size cannot be negative")
			} else {
				require.NoError(t, err)
			}
			require.Equal(t, testCase.expectedResult, formattedSize, "sizes do not match")
		})
	}
}
