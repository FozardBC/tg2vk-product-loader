package config

import (
	"log"
	"os"
	"path/filepath"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env           string `yaml:"env" env-default:"local" env-required:"true"`
	TelegramToken string `yaml:"telegram_token" env-required:"true"`
	VkToken       string `yaml:"vk_token" env-required:"true"`
	VkGroupId     int    `yaml:"VK_GROUP_ID" env-required:"true"`
}

func MustLoad() *Config {
	projectDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("can't get projDirectory to set confing file: %s", err.Error())
	}

	configPath := filepath.Join(projectDir, "..", "..", "config", "cfg.yaml")

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s", configPath)
	}

	var cfg Config

	if err = cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("can't read config:%s", err.Error())
	}

	return &cfg

}
