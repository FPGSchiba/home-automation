package util

import (
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	Database DatabaseConfig `yaml:"database"`
	Security SecurityConfig `yaml:"security"`
	TLS      TLSConfig      `yaml:"tls"`
}

type DatabaseConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
}

type SecurityConfig struct {
	TokenExpiration int `yaml:"tokenExpiration"`
}

type TLSConfig struct {
	Enabled bool   `yaml:"enabled"`
	Cert    string `yaml:"cert"`
	Key     string `yaml:"key"`
	Chain   string `yaml:"chain"`
}

var configFilepath string

func SetConfigFilePath(path string) {
	configFilepath = path
}

func getFilePath() string {
	if configFilepath == "" {
		println("config file path is empty")
		configFilepath = "config.yaml"
	}
	return configFilepath
}

func (c *Config) GetConfig() *Config {
	yamlFile, err := os.ReadFile(getFilePath())
	if err != nil {
		log.WithFields(log.Fields{
			"error":     err,
			"component": "config",
			"func":      "getConfig",
			"file":      configFilepath,
		}).Error("Error reading config file")
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	return c
}
