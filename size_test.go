package code

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func createTestFile(t *testing.T, testFilePath string, size int) {
	t.Helper()
	bytes := strings.Repeat("1", size)
	err := os.WriteFile(testFilePath, []byte(bytes), 0644)
	require.NoError(t, err, "unable to create test file")
}

func setupRecursionTest(t *testing.T) string {
	dir := t.TempDir()
	deepDirPath := filepath.Join(dir, "level1", "level2", "level3")
	err := os.MkdirAll(deepDirPath, 0755)
	require.NoError(t, err)

	file20B := filepath.Join(dir, "file20")
	createTestFile(t, file20B, 20)

	hFileLvl1 := filepath.Join(dir, "level1", ".file_level1_30")
	createTestFile(t, hFileLvl1, 30)

	fileLvl2 := filepath.Join(dir, "level1", "level2", "file_level2_2080")
	createTestFile(t, fileLvl2, 2080)

	fileLvl3 := filepath.Join(dir, "level1", "level2", "level3", "file_level2_1700")
	createTestFile(t, fileLvl3, 1700)
	return dir
}

func TestGetPathSize(t *testing.T) {
	tempDir := t.TempDir()

	file5B := filepath.Join(tempDir, "file5.txt")
	createTestFile(t, file5B, 5)

	file10B := filepath.Join(tempDir, "file10.txt")
	createTestFile(t, file10B, 10)

	hiddenFile25 := filepath.Join(tempDir, ".hidden25")
	createTestFile(t, hiddenFile25, 25)

	tempDir2 := t.TempDir()
	file1038B := filepath.Join(tempDir2, ".hidden1038")
	createTestFile(t, file1038B, 1038)

	dir := t.TempDir()
	hiddenDir, err := os.MkdirTemp(dir, ".hidden")
	require.NoError(t, err)

	file2 := filepath.Join(hiddenDir, "file2")
	createTestFile(t, file2, 2)

	recursiveDir := setupRecursionTest(t)

	testCases := []struct {
		name         string
		path         string
		expectedSize int64
		withHidden   bool
		recursive    bool
		expectError  bool
	}{
		{
			name:         "1 file of 5 bytes",
			path:         file5B,
			expectedSize: 5,
			withHidden:   false,
			recursive:    false,
			expectError:  false,
		},
		{
			name:         "1 file of 10 bytes",
			path:         file10B,
			expectedSize: 10,
			withHidden:   false,
			recursive:    false,
			expectError:  false,
		},
		{
			name:         "1 directory with 2 files",
			path:         tempDir,
			expectedSize: 15,
			withHidden:   false,
			recursive:    false,
			expectError:  false,
		},
		{
			name:         "Empty path returns error",
			path:         "",
			expectedSize: 0,
			withHidden:   false,
			recursive:    false,
			expectError:  true,
		},
		{
			name:         "1 directory with 2 regular files and 1 hidden: show hidden",
			path:         tempDir,
			expectedSize: 40,
			withHidden:   true,
			recursive:    false,
			expectError:  false,
		},
		{
			name:         "1 directory with 2 regular files and 1 hidden: do not show hidden",
			path:         tempDir,
			expectedSize: 15,
			withHidden:   false,
			recursive:    false,
			expectError:  false,
		},
		{
			name:         "1 directory with 1 hidden file: show hidden",
			path:         tempDir2,
			expectedSize: 1038,
			withHidden:   true,
			recursive:    false,
			expectError:  false,
		},
		{
			name:         "1 directory with 1 hidden file: do not show hidden",
			path:         tempDir2,
			expectedSize: 0,
			withHidden:   false,
			recursive:    false,
			expectError:  false,
		},
		{
			name:         "1 hidden directory with 1 regular file: do not show hidden",
			path:         hiddenDir,
			expectedSize: 0,
			withHidden:   false,
			recursive:    false,
			expectError:  false,
		},
		{
			name:         "1 dir with two nested levels",
			path:         recursiveDir,
			expectedSize: 3800,
			withHidden:   false,
			recursive:    true,
			expectError:  false,
		},
		{
			name:         "1 dir with two nested levels: 1 level only",
			path:         recursiveDir,
			expectedSize: 20,
			withHidden:   false,
			recursive:    false,
			expectError:  false,
		},
		{
			name:         "1 dir with 1 nested levels: with hidden",
			path:         recursiveDir,
			expectedSize: 3830,
			withHidden:   true,
			recursive:    true,
			expectError:  false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			size, err := GetSize(tc.path, tc.withHidden, tc.recursive)
			if tc.expectError {
				require.ErrorIs(t, err, ErrorEmptyPath)
			} else {
				require.NoError(t, err)
			}
			require.Equal(t, tc.expectedSize, size)
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
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			formattedSize, err := FormatSize(int64(tc.size), tc.humanReadable)
			if tc.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			require.Equal(t, tc.expectedResult, formattedSize)
		})
	}
}
