package config

import (
	"flag"
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env          string `yaml:"env" env:"ENV" env-default:"production"`
	Storage_path string `yaml:"storage_path" env:"STORAGE_PATH" env-default:"storage/storage.db"`
	Http_server  `yaml:"http_server"`
}

type Http_server struct {
	Address string `yaml:"address" env:"ADDRESS" env-default:"0.0.0.0:3000"`
}

func MustLoad() *Config {
	var configPath string
	configPath = os.Getenv("CONFIG_PATH")

	if configPath == "" {
		flags := flag.String("config", "", "Path to config file")
		flag.Parse()

		configPath = *flags

		if configPath == "" {
			log.Fatal("Config path is not set")
		}
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatal("Cofig file does not exist", configPath)
	}

	var cfg Config

	err := cleanenv.ReadConfig(configPath, &cfg)

	if err != nil {
		log.Fatal("Cannot read config file", err)
	}

	return &cfg

}
