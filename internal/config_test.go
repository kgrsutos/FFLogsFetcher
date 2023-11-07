package internal

import (
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func setupTestFile(content string) (string, func(), error) {
	tmpfile, err := os.CreateTemp("", "testconfig*.yaml")
	if err != nil {
		return "", nil, err
	}
	_, err = tmpfile.WriteString(content)
	if err != nil {
		return "", nil, err
	}
	cleanup := func() {
		os.Remove(tmpfile.Name())
	}
	return tmpfile.Name(), cleanup, nil
}

func TestLoadRequestConfig_OK(t *testing.T) {
	content := `
reportIDs:
  - "ABC123"
  - "XYZ789"
`
	path, cleanup, err := setupTestFile(content)
	if err != nil {
		t.Fatal(err)
	}
	defer cleanup()
	config, err := LoadRequestConfig(path)
	if err != nil {
		t.Fatal(err)
	}
	want := RequestConfig{
		ReportIDs: []string{"ABC123", "XYZ789"},
	}
	if diff := cmp.Diff(want, config); diff != "" {
		t.Errorf("LoadConfig() mismatch (-want +got):\n%s", diff)
	}
}

func TestLoadRequestConfig_FileError(t *testing.T) {
	path := "/path/to/nonexistent/file.yaml"
	_, err := LoadRequestConfig(path)
	if err == nil {
		t.Error("expected an error but got none")
	}
}

func TestLoadRequestConfig_UnmarshalError(t *testing.T) {
	content := `
reportIDs: "not a list"
`
	path, cleanup, err := setupTestFile(content)
	if err != nil {
		t.Fatal(err)
	}
	defer cleanup()
	_, err = LoadRequestConfig(path)
	if err == nil {
		t.Error("expected an error but got none")
	}
}
