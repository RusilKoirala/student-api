package config

import (
	"flag"
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	// Env default is production
	Env          string `yaml:"env" env:"ENV" env-required:"true" env-default:"production"`
	Storage_path string `yaml:"storage_path" env-required:"true"`
	Http_server  `yaml:"http_server"`
}

type Http_server struct {
	Address string
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
