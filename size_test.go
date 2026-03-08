package code

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

type testCase struct {
	name         string
	path         string
	expectedSize int64
	withHidden   bool
	recursive    bool
	expectError  bool
}

func TestGetSize_Basic(t *testing.T) {
	testDir1 := filepath.Join("testdata", "dir1")
	testCases := []testCase{
		{
			name:         "5B file",
			path:         filepath.Join(testDir1, "file5"),
			expectedSize: 5,
			withHidden:   false,
			recursive:    false,
			expectError:  false,
		},
		{
			name:         "10B file",
			path:         filepath.Join(testDir1, "file10"),
			expectedSize: 10,
			withHidden:   false,
			recursive:    false,
			expectError:  false,
		},
		{
			name:         "dir with files",
			path:         testDir1,
			expectedSize: 15,
			withHidden:   false,
			recursive:    false,
			expectError:  false,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			size, err := GetSize(tc.path, tc.withHidden, tc.recursive)
			require.NoError(t, err)
			require.Equal(t, tc.expectedSize, size)
		})
	}
}

func TestGetSize_EmptyPath(t *testing.T) {
	_, err := GetSize("", false, false)
	require.ErrorIs(t, err, ErrorEmptyPath)
}

func TestGetSize_HiddenFiles(t *testing.T) {
	testDir1 := filepath.Join("testdata", "dir1")
	testCases := []testCase{
		{
			name:         "show hidden",
			path:         testDir1,
			expectedSize: 40,
			withHidden:   true,
			recursive:    false,
			expectError:  false,
		},
		{
			name:         "ignore hidden",
			path:         testDir1,
			expectedSize: 15,
			withHidden:   false,
			recursive:    false,
			expectError:  false,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			size, err := GetSize(tc.path, tc.withHidden, tc.recursive)
			require.NoError(t, err)
			require.Equal(t, tc.expectedSize, size)
		})
	}
}

func TestGetSize_SingleFile(t *testing.T) {
	testDir2 := filepath.Join("testdata", "dir2")
	testCases := []testCase{
		{
			name:         "show hidden",
			path:         testDir2,
			expectedSize: 1038,
			withHidden:   true,
			recursive:    false,
			expectError:  false,
		},
		{
			name:         "ignore hidden",
			path:         testDir2,
			expectedSize: 0,
			withHidden:   false,
			recursive:    false,
			expectError:  false,
		},
		{
			name:         "ignore hidden dir",
			path:         filepath.Join("testdata", ".hidden"),
			expectedSize: 0,
			withHidden:   false,
			recursive:    false,
			expectError:  false,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			size, err := GetSize(tc.path, tc.withHidden, tc.recursive)
			require.NoError(t, err)
			require.Equal(t, tc.expectedSize, size)
		})
	}
}

func TestGetSize_Recursive(t *testing.T) {
	recDir := filepath.Join("testdata", "recDir")
	testCases := []testCase{
		{
			name:         "with recursion",
			path:         recDir,
			expectedSize: 3800,
			withHidden:   false,
			recursive:    true,
			expectError:  false,
		},
		{
			name:         "1 level only",
			path:         recDir,
			expectedSize: 20,
			withHidden:   false,
			recursive:    false,
			expectError:  false,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			size, err := GetSize(tc.path, tc.withHidden, tc.recursive)
			require.NoError(t, err)
			require.Equal(t, tc.expectedSize, size)
		})
	}
}

func TestGetSize_Recursive_All(t *testing.T) {
	path := filepath.Join("testdata", "recDir")
	size, err := GetSize(path, true, true)
	require.NoError(t, err)
	require.Equal(t, int64(3830), size)
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
			name:           "No format",
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
