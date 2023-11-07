package query

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
	// テスト用のYAMLファイルをセットアップ
	path, cleanup, err := setupTestFile(content)
	require.NoError(t, err)
	defer cleanup()
	// LoadConfig関数をテスト
	config, err := LoadConfig(path)
	require.NoError(t, err)
	assert.Equal(t, []string{"ABC123", "XYZ789"}, config.ReportIDs)
}

func TestLoadConfig_FileError(t *testing.T) {
	path := "/path/to/nonexistent/file.yaml"
	_, err := LoadConfig(path)
	assert.Error(t, err)
}

func TestLoadConfig_UnmarshalError(t *testing.T) {
	content := `
reportIDs: "not a list"
`
	path, cleanup, err := setupTestFile(content)
	require.NoError(t, err)
	defer cleanup()
	_, err = LoadConfig(path)
	assert.Error(t, err)
}
