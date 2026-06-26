package config

import (
	"fmt"
	"log"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Postgres Postgres `yaml:"postgres"`
	Bot      Kafka    `yaml:"kafka"`
}

type Postgres struct {
	Host     string `yaml:"host"`
	Port     int64  `yaml:"port"`
	Schema   string `yaml:"schema"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
}

type Kafka struct {
	Host          string `yaml:"host"`
	Port          int64  `yaml:"port"`
	Topic         string `yaml:"topic"`
	ConsumerGroup string `yaml:"consumer_group"`
}

func (c *Config) GetPostgresURL() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		c.Postgres.Username, c.Postgres.Password, c.Postgres.Host, c.Postgres.Port, c.Postgres.Database)
}

func GetConfig(configPath string) *Config {
	if configPath == "" {
		log.Fatal("config path is required")
	}

	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatal(err.Error())
	}

	return &cfg
}
