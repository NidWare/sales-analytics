package config

import (
	_ "fmt"
	_ "log"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	BotToken   string  `yaml:"bot_token"`
	AsanaToken string  `yaml:"asana_token"`
	AdminIDs   []int64 `yaml:"admin_ids"`
}

func LoadConfig(filename string) (*Config, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var config Config
	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
