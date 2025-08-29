package pkgmgr

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/project-copacetic/copacetic/pkg/buildkit"
	"github.com/project-copacetic/copacetic/pkg/types/unversioned"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestGetPackageManager tests the GetPackageManager function.
func TestGetPackageManager(t *testing.T) {
	// Create a mock config and workingFolder
	config := &buildkit.Config{}
	workingFolder := "/tmp"

	t.Run("should return an apkManager for alpine", func(t *testing.T) {
		// Call the GetPackageManager function with "alpine" as osType
		manager, err := GetPackageManager("alpine", "1.0", config, workingFolder)

		// Assert that there is no error and the manager is not nil
		assert.NoError(t, err)
		assert.NotNil(t, manager)

		// Assert that the manager is an instance of apkManager
		assert.IsType(t, &apkManager{}, manager)
	})

	t.Run("should return a dpkgManager for debian", func(t *testing.T) {
		// Call the GetPackageManager function with "debian" as osType
		manager, err := GetPackageManager("debian", "1.0", config, workingFolder)

		// Assert that there is no error and the manager is not nil
		assert.NoError(t, err)
		assert.NotNil(t, manager)

		// Assert that the manager is an instance of dpkgManager
		assert.IsType(t, &dpkgManager{}, manager)
	})

	t.Run("should return a dpkgManager for ubuntu", func(t *testing.T) {
		// Call the GetPackageManager function with "ubuntu" as osType
		manager, err := GetPackageManager("ubuntu", "1.0", config, workingFolder)

		// Assert that there is no error and the manager is not nil
		assert.NoError(t, err)
		assert.NotNil(t, manager)

		// Assert that the manager is an instance of dpkgManager
		assert.IsType(t, &dpkgManager{}, manager)
	})

	t.Run("should return an rpmManager for cbl-mariner", func(t *testing.T) {
		// Call the GetPackageManager function with "cbl-mariner" as osType
		manager, err := GetPackageManager("cbl-mariner", "1.0", config, workingFolder)

		// Assert that there is no error and the manager is not nil
		assert.NoError(t, err)
		assert.NotNil(t, manager)

		// Assert that the manager is an instance of rpmManager
		assert.IsType(t, &rpmManager{}, manager)
	})

	t.Run("should return an rpmManager for azurelinux", func(t *testing.T) {
		// Call the GetPackageManager function with "azurelinux" as osType
		manager, err := GetPackageManager("azurelinux", "1.0", config, workingFolder)

		// Assert that there is no error and the manager is not nil
		assert.NoError(t, err)
		assert.NotNil(t, manager)

		// Assert that the manager is an instance of rpmManager
		assert.IsType(t, &rpmManager{}, manager)
	})

	t.Run("should return an rpmManager for redhat", func(t *testing.T) {
		// Call the GetPackageManager function with "redhat" as osType
		manager, err := GetPackageManager("redhat", "1.0", config, workingFolder)

		// Assert that there is no error and the manager is not nil
		assert.NoError(t, err)
		assert.NotNil(t, manager)

		// Assert that the manager is an instance of rpmManager
		assert.IsType(t, &rpmManager{}, manager)
	})

	t.Run("should return an error for unsupported osType", func(t *testing.T) {
		// Call the GetPackageManager function with "unsupported" as osType
		manager, err := GetPackageManager("unsupported", "", config, workingFolder)

		// Assert that there is an error and the manager is nil
		assert.Error(t, err)
		assert.Nil(t, manager)
	})
}

func IsValid(version string) bool {
	return version != "invalid"
}

func LessThan(v1, v2 string) bool {
	// Simplistic comparison for testing
	return v1 < v2
}

func TestGetUniqueLatestUpdates(t *testing.T) {
	cmp := VersionComparer{IsValid, LessThan}

	tests := []struct {
		name          string
		updates       unversioned.UpdatePackages
		ignoreErrors  bool
		want          unversioned.UpdatePackages
		expectedError string
	}{
		{
			name:          "empty updates",
			updates:       unversioned.UpdatePackages{},
			ignoreErrors:  false,
			want:          nil,
			expectedError: "no patchable vulnerabilities found",
		},
		{
			name: "valid updates",
			updates: unversioned.UpdatePackages{
				{Name: "pkg1", FixedVersion: "1.0"},
				{Name: "pkg1", FixedVersion: "2.0"},
			},
			ignoreErrors: false,
			want: unversioned.UpdatePackages{
				{Name: "pkg1", FixedVersion: "2.0"},
			},
			expectedError: "",
		},
		{
			name: "updates with invalid version",
			updates: unversioned.UpdatePackages{
				{Name: "pkg1", FixedVersion: "invalid"},
			},
			ignoreErrors:  false,
			want:          nil,
			expectedError: "invalid version invalid found for package pkg1",
		},
		{
			name: "ignore errors",
			updates: unversioned.UpdatePackages{
				{Name: "pkg1", FixedVersion: "invalid"},
			},
			ignoreErrors:  true,
			want:          unversioned.UpdatePackages{},
			expectedError: "",
		},
		{
			name: "Updates with the same highest version",
			updates: unversioned.UpdatePackages{
				{Name: "pkg2", FixedVersion: "2.0"},
				{Name: "pkg1", FixedVersion: "1.0"},
				{Name: "pkg2", FixedVersion: "2.0"},
				{Name: "pkg1", FixedVersion: "1.0"},
			},
			ignoreErrors: false,
			want: unversioned.UpdatePackages{
				{Name: "pkg1", FixedVersion: "1.0"},
				{Name: "pkg2", FixedVersion: "2.0"},
			},
			expectedError: "",
		},
		{
			name: "Invalid versions with ignoreErrors true",
			updates: unversioned.UpdatePackages{
				{Name: "pkg1", FixedVersion: "invalid"},
				{Name: "pkg2", FixedVersion: "3.0"},
				{Name: "pkg3", FixedVersion: "invalid"},
			},
			ignoreErrors: true,
			want: unversioned.UpdatePackages{
				{Name: "pkg2", FixedVersion: "3.0"},
			},
			expectedError: "",
		},
		{
			name: "Updates with decreasing versions",
			updates: unversioned.UpdatePackages{
				{Name: "pkg1", FixedVersion: "2.0"},
				{Name: "pkg1", FixedVersion: "1.5"},
				{Name: "pkg1", FixedVersion: "3.0"},
			},
			ignoreErrors: false,
			want: unversioned.UpdatePackages{
				{Name: "pkg1", FixedVersion: "3.0"},
			},
			expectedError: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetUniqueLatestUpdates(tt.updates, cmp, tt.ignoreErrors)
			if err != nil {
				if tt.expectedError == "" {
					t.Errorf("GetUniqueLatestUpdates() unexpected error = %v", err)
				} else if !strings.Contains(err.Error(), tt.expectedError) {
					t.Errorf("GetUniqueLatestUpdates() error = %v, wantErrMsg %v", err, tt.expectedError)
				}
			} else if tt.expectedError != "" {
				t.Errorf("GetUniqueLatestUpdates() expected error %v, got none", tt.expectedError)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("%s: got = %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}

// Mock PackageInfoReader for testing
type mockPackageInfoReader struct {
	nameMap    map[string]string
	versionMap map[string]string
	errors     map[string]error
}

func (m *mockPackageInfoReader) GetName(filename string) (string, error) {
	if err, ok := m.errors[filename]; ok {
		return "", err
	}
	if name, ok := m.nameMap[filename]; ok {
		return name, nil
	}
	return "", fmt.Errorf("no name mapping for %s", filename)
}

func (m *mockPackageInfoReader) GetVersion(filename string) (string, error) {
	if err, ok := m.errors[filename]; ok {
		return "", err
	}
	if version, ok := m.versionMap[filename]; ok {
		return version, nil
	}
	return "", fmt.Errorf("no version mapping for %s", filename)
}

func TestGetValidatedUpdatesMap(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "pkgmgr-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	cmp := VersionComparer{IsValid, LessThan}

	tests := []struct {
		name           string
		updates        unversioned.UpdatePackages
		files          []string
		reader         *mockPackageInfoReader
		expectedResult UpdateMap
		expectedError  string
	}{
		{
			name:           "empty staging directory",
			updates:        unversioned.UpdatePackages{{Name: "pkg1", FixedVersion: "1.0"}},
			files:          []string{}, // No files in staging
			reader:         &mockPackageInfoReader{},
			expectedResult: nil,
			expectedError:  "",
		},
		{
			name:    "successful validation with matching packages",
			updates: unversioned.UpdatePackages{{Name: "pkg1", FixedVersion: "2.0"}},
			files:   []string{"pkg1_2.0.deb"},
			reader: &mockPackageInfoReader{
				nameMap:    map[string]string{"pkg1_2.0.deb": "pkg1"},
				versionMap: map[string]string{"pkg1_2.0.deb": "2.0"},
			},
			expectedResult: UpdateMap{"pkg1": {Version: "2.0", Filename: "pkg1_2.0.deb"}},
			expectedError:  "",
		},
		{
			name:    "version mismatch error",
			updates: unversioned.UpdatePackages{{Name: "pkg1", FixedVersion: "2.0"}},
			files:   []string{"pkg1_1.5.deb"},
			reader: &mockPackageInfoReader{
				nameMap:    map[string]string{"pkg1_1.5.deb": "pkg1"},
				versionMap: map[string]string{"pkg1_1.5.deb": "1.5"},
			},
			expectedResult: nil,
			expectedError:  "downloaded package pkg1 version 1.5 lower than required 2.0",
		},
		{
			name:    "invalid version in downloaded package",
			updates: unversioned.UpdatePackages{{Name: "pkg1", FixedVersion: "2.0"}},
			files:   []string{"pkg1_invalid.deb"},
			reader: &mockPackageInfoReader{
				nameMap:    map[string]string{"pkg1_invalid.deb": "pkg1"},
				versionMap: map[string]string{"pkg1_invalid.deb": "invalid"},
			},
			expectedResult: nil,
			expectedError:  "invalid version invalid found for package pkg1",
		},
		{
			name:    "package name parsing error",
			updates: unversioned.UpdatePackages{{Name: "pkg1", FixedVersion: "1.0"}},
			files:   []string{"broken.deb"},
			reader: &mockPackageInfoReader{
				errors: map[string]error{"broken.deb": fmt.Errorf("parse error")},
			},
			expectedResult: nil,
			expectedError:  "parse error",
		},
		{
			name:    "package version parsing error",
			updates: unversioned.UpdatePackages{{Name: "pkg1", FixedVersion: "1.0"}},
			files:   []string{"pkg1.deb"},
			reader: &mockPackageInfoReader{
				nameMap: map[string]string{"pkg1.deb": "pkg1"},
				errors:  map[string]error{"pkg1.deb": fmt.Errorf("version parse error")},
			},
			expectedResult: nil,
			expectedError:  "version parse error",
		},
		{
			name:    "unexpected package not in updates - should be ignored with warning",
			updates: unversioned.UpdatePackages{{Name: "pkg1", FixedVersion: "1.0"}},
			files:   []string{"pkg2_1.0.deb"},
			reader: &mockPackageInfoReader{
				nameMap:    map[string]string{"pkg2_1.0.deb": "pkg2"},
				versionMap: map[string]string{"pkg2_1.0.deb": "1.0"},
			},
			expectedResult: UpdateMap{"pkg1": {Version: "1.0", Filename: ""}}, // pkg1 exists but no matching file
			expectedError:  "",
		},
		{
			name: "multiple packages successful validation",
			updates: unversioned.UpdatePackages{
				{Name: "pkg1", FixedVersion: "1.0"},
				{Name: "pkg2", FixedVersion: "2.0"},
			},
			files: []string{"pkg1_1.0.deb", "pkg2_2.0.deb"},
			reader: &mockPackageInfoReader{
				nameMap:    map[string]string{"pkg1_1.0.deb": "pkg1", "pkg2_2.0.deb": "pkg2"},
				versionMap: map[string]string{"pkg1_1.0.deb": "1.0", "pkg2_2.0.deb": "2.0"},
			},
			expectedResult: UpdateMap{
				"pkg1": {Version: "1.0", Filename: "pkg1_1.0.deb"},
				"pkg2": {Version: "2.0", Filename: "pkg2_2.0.deb"},
			},
			expectedError: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test files in the temporary directory
			for _, filename := range tt.files {
				filePath := filepath.Join(tempDir, filename)
				err := os.WriteFile(filePath, []byte("dummy content"), 0o644)
				require.NoError(t, err)
			}

			// Call the function under test
			result, err := GetValidatedUpdatesMap(tt.updates, cmp, tt.reader, tempDir)

			// Check error expectations
			if tt.expectedError != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				if tt.expectedResult == nil {
					assert.Nil(t, result)
				} else {
					assert.Equal(t, tt.expectedResult, result)
				}
			}

			// Clean up test files for next iteration
			for _, filename := range tt.files {
				filePath := filepath.Join(tempDir, filename)
				os.Remove(filePath)
			}
		})
	}
}

func TestGetValidatedUpdatesMap_DirectoryErrors(t *testing.T) {
	cmp := VersionComparer{IsValid, LessThan}
	reader := &mockPackageInfoReader{}
	updates := unversioned.UpdatePackages{{Name: "pkg1", FixedVersion: "1.0"}}

	// Test with non-existent directory
	result, err := GetValidatedUpdatesMap(updates, cmp, reader, "/non/existent/path")
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "no such file or directory")
}
