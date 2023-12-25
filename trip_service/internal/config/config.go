package config

import (
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	Http  Http  `yaml:"http"`
	DB    DB    `yaml:"db"`
	Kafka Kafka `yaml:"kafka"`
}

type Http struct {
	Port int `yaml:"port"`
}
type DB struct {
	DSN           string `yaml:"dsn"`
	MigrationsDir string `yaml:"migrations_dir"`
}
type Kafka struct {
	Address string `yaml:"address"`
}

func NewConfig(filePath string) (*Config, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	var cfg Config
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}
