package config

import (
	"flag"
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type HTTPServerConfig struct {
	Address string `yaml:"address"`
}

type Config struct {
	Env         string           `yaml:"env"`
	Version     string           `yaml:"version"`
	Description string           `yaml:"description"`
	StoragePath string           `yaml:"storage_path"`
	HTTPServer  HTTPServerConfig `yaml:"http_server"`
}

func MustLoadConfig() *Config {
	var configPath string

	configPath = os.Getenv("PANES_CONFIG_PATH")
	if configPath == "" {
		flags := flag.String("config", "config/panes.yaml", "Path to the configuration file")
		flag.Parse()
		configPath = *flags
		if configPath == "" {
			log.Fatal("Configuration file path is not set. Please set the PANES_CONIG_PATH environment variable or use the --config flag.")
		}
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("Configuration file does not exist at path: %s", configPath)
	}
	var config Config
	err := cleanenv.ReadConfig(configPath, &config)

	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	return &config

}
