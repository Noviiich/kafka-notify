package config

import (
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load("C:/Users/nowik/VSCode/Go/kafka-notify/.env")
	if err != nil {
		log.Print("Error loading .env file: ", err)
	}
}

type Config struct {
	Env         string `yaml:"env" env-default:"local"`
	StoragePath string `yaml:"storage_path" env-required:"true"`
	// Server      Server      `yaml:"server:"`
	KafkaConfig KafkaConfig `yaml:"kafka"`
	Producer    `yaml:"producer"`
	Consumer    `yaml:"consumer"`
}

// type Server struct {
// 	Producer Producer `yaml:"producer"`
// 	Consumer Consumer `yaml:"consumer"`
// }

type KafkaConfig struct {
	Brokers       []string `yaml:"brokers"`
	Topic         string   `yaml:"topic"`
	ConsumerGroup string   `yaml:"consumer_group"`
}

type Producer struct {
	Address string `yaml:"address"`
}

type Consumer struct {
	Address string `yaml:"address"`
}

func MustLoad() *Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH is not set")
	}

	// check if file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s", configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}

	return &cfg
}
