package query

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	ReportIDs []string `yaml:"reportIDs"`
}

func LoadConfig(path string) (Config, error) {
	var config Config
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
