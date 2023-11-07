package query

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

func TestLoadConfig_OK(t *testing.T) {
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
	config, err := LoadConfig(path)
	if err != nil {
		t.Fatal(err)
	}
	want := Config{
		ReportIDs: []string{"ABC123", "XYZ789"},
	}
	if diff := cmp.Diff(want, config); diff != "" {
		t.Errorf("LoadConfig() mismatch (-want +got):\n%s", diff)
	}
}

func TestLoadConfig_FileError(t *testing.T) {
	path := "/path/to/nonexistent/file.yaml"
	_, err := LoadConfig(path)
	if err == nil {
		t.Error("expected an error but got none")
	}
}

func TestLoadConfig_UnmarshalError(t *testing.T) {
	content := `
reportIDs: "not a list"
`
	path, cleanup, err := setupTestFile(content)
	if err != nil {
		t.Fatal(err)
	}
	defer cleanup()
	_, err = LoadConfig(path)
	if err == nil {
		t.Error("expected an error but got none")
	}
}
