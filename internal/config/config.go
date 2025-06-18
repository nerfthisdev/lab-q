package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Server struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

type Bot struct {
	APIToken   string `yaml:"token"`
	WebHookURL string `yaml:"webHookURL"`
}

type Config struct {
	Server Server `yaml:"server"`
	Bot    Bot    `yaml:"bot"`
}

func GetConfiguration(configPath string, cfg interface{}) error {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		return nil
	}
	return nil
}
