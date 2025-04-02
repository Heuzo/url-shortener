package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
	"time"
)

type Config struct {
	Env         string `yaml:"env" env:"ENV" env-required:"true"`
	StoragePath string `yaml:"storage_path" env:"STORAGE_PATH" env-required:"true"`
	AliasLenght int    `yaml:"alias_lenght" env:"ALIAS_LENGTH" env-default:"6"`
	HTTPServer  `yaml:"http_server"`
}

type HTTPServer struct {
	Address     string        `yaml:"address" env:"ADDRESS" env-required:"true"`
	Timeout     time.Duration `yaml:"timeout" env:"TIMEOUT" env-required:"true"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env:"IDLE_TIMEOUT" env-required:"true"`
	User        string        `yaml:"user" env:"HTTP_SERVER_USER" env-required:"true"`
	Password    string        `yaml:"password" env:"HTTP_SERVER_PASSWORD" env-required:"true"`
}

func MustLoad() *Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH environment variable not set")
	}
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("does not exist config file by path: %s", configPath)
	}
	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}
	return &cfg
}
