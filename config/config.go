package config

import (
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
)

type Config struct {
	App struct {
		Port        string `yaml:"port"`
		Host        string `yaml:"host"`
		LogLevel    string `yaml:"log-level"`
	} `yaml:"app"`
}

func Setup(configPath string) (config Config, err error) {
	filename, _ := filepath.Abs(configPath)
	yamlFile, err := os.ReadFile(filename)
	if err != nil {
		return config, err
	}
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		return
	}

	if config.App.LogLevel == "" {
		config.App.LogLevel = "Error"
	}

	return
}
