package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Port                string   `yaml:"port"`
	BackendServers      []string `yaml:"backend_servers"`
	HealthCheckInterval int      `yaml:"health_check_interval"`
	Algorithm           string   `yaml:"algorithm"`
}

func LoadConfig(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
