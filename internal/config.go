package internal

import (
	"os"

	"gopkg.in/yaml.v3"
)

type RequestConfig struct {
	ReportIDs []string `yaml:"reportIDs"`
}

func LoadRequestConfig(path string) (RequestConfig, error) {
	var config RequestConfig
	data, err := os.ReadFile(path)
	if err != nil {
		return config, err
	}
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return config, err
	}
	return config, nil
}
