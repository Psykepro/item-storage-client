package config

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

const (
	Path = "./config/config.yml"
)

type Config struct {
	RabbitMQ *RabbitMQ `yaml:"rabbit_mq"`
	Logger   *Logger   `yaml:"logging"`
}

type RabbitMQ struct {
	Host      string    `yaml:"host"`
	Port      string    `yaml:"port"`
	User      string    `yaml:"user"`
	Password  string    `yaml:"password"`
	Queue     Queue     `yaml:"queue"`
	Publisher Publisher `yaml:"publisher"`
}

type Queue struct {
	Name       string         `yaml:"name"`
	Durable    bool           `yaml:"durable"`
	AutoDelete bool           `yaml:"auto_delete"`
	Exclusive  bool           `yaml:"exclusive"`
	NoWait     bool           `yaml:"no_wait"`
	Args       map[string]any `yaml:"args"`
}

type Publisher struct {
	Exchange    string `yaml:"exchange"`
	Mandatory   bool   `yaml:"mandatory"`
	Immediate   bool   `yaml:"immediate"`
	ContentType string `yaml:"content_type"`
}

type Logger struct {
	Level    string `yaml:"level"`
	Mode     string `yaml:"mode"`
	Encoding string `yaml:"encoding"`
}

func GetConfig(path string) (*Config, error) {
	var config Config
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, err
	}

	filename, err := filepath.Abs(path)
	if err != nil {
		log.Fatal(err)
	}

	yamlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}

	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		log.Fatal(err)
	}

	return &config, err
}
