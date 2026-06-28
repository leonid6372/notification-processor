package config

import (
	"fmt"
	"log"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Scheduler *Scheduler `yaml:"scheduler"`
	Sender    *Sender    `yaml:"sender"`
	Postgres  *Postgres  `yaml:"postgres"`
	Kafka     *Kafka     `yaml:"kafka"`
}

type Scheduler struct {
	BatchSize    int           `yaml:"batch_size"` // get X notifications from DB to process
	WorkersCount int           `yaml:"workers_count"`
	TaskTimeout  time.Duration `yaml:"task_timeout"` // task timeout to avoid zombie state in DB
}

type Sender struct {
	RetryCount int           `yaml:"retry_count"` // try to send notification X times before mark as failed
	MinDelay   time.Duration `yaml:"min_delay"`
}

type Postgres struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Schema   string `yaml:"schema"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
}

func (c *Config) GetPostgresURL() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		c.Postgres.Username, c.Postgres.Password, c.Postgres.Host, c.Postgres.Port, c.Postgres.Database)
}

type Kafka struct {
	Host      string `yaml:"host"`
	Port      int    `yaml:"port"`
	Topic     string `yaml:"topic"`
	GroupID   string `yaml:"group_id"`
	BatchSize int    `yaml:"batch_size"` // commit messages to kafka every X messages
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
