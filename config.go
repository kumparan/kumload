package kumload

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

const (
	// LogLevelDebug :nodoc:
	LogLevelDebug = "debug"
)

// Batch :nodoc:
type Batch struct {
	Size     int64 `yaml:"size"`
	Interval int64 `yaml:"interval"`
}

// DatabaseDetail :nodoc:
type DatabaseDetail struct {
	Name       string `yaml:"name"`
	Table      string `yaml:"table"`
	PrimaryKey string `yaml:"primary_key"`
	OrderKey   string `yaml:"order_key"`
}

// Database :nodoc:
type Database struct {
	Host     string         `yaml:"host"`
	Username string         `yaml:"username"`
	Password string         `yaml:"password"`
	Source   DatabaseDetail `yaml:"source"`
	Target   DatabaseDetail `yaml:"target"`
}

// Config :nodoc:
type Config struct {
	LogLevel string            `yaml:"log_level"`
	Batch    Batch             `yaml:"batch"`
	Database Database          `yaml:"database"`
	Mappings map[string]string `yaml:"mappings"`
	Script   string
}

// ParseConfig parse config from given path
func ParseConfig(path string) (*Config, error) {
	configByte, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config Config
	err = yaml.Unmarshal(configByte, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
