package config

import (
	"gopkg.in/yaml.v2"
	"os"
)

type Config struct {
	BotToken   string  `yaml:"bot_token"`
	AsanaToken string  `yaml:"asana_token"`
	AdminIDs   []int64 `yaml:"admin_ids"`
}

func LoadConfig(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var cfg Config
	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(&cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
