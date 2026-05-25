package config

import (
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Server   ServerConfig `yaml:"http_server"`
	Database DBConfig     `yaml:"database"`
	Logger   LogConfig    `yaml:"logger"`
}

type DBConfig struct {
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	DBName   string `yaml:"name"`
	Reload   bool   `yaml:"reload"`
}

type ServerConfig struct {
	Address string `yaml:"address"`
}

type LogConfig struct {
	LogFile  string `yaml:"logfile"`
	LogLevel string `yaml:"loglevel"`
}

func New() *Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "./config/local.yaml"
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s", configPath)
	}

	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("can't read config: %s", err)
	}

	return &cfg
}
