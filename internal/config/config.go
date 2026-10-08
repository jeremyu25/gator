package config

import (
	"encoding/json"
	"fmt"
	"os"
)

const configFileName = "/.gatorconfig.json"

type Config struct {
	DbUrl           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func Read() Config {
	configPath, err := getConfigPath()
	if err != nil {
		return Config{}
	}
	jsonFile, err := os.Open(configPath)
	if err != nil {
		return Config{}
	}
	defer jsonFile.Close()
	decoder := json.NewDecoder(jsonFile)
	var cfg Config
	if err = decoder.Decode(&cfg); err != nil {
		return Config{}
	}
	return cfg
}

func getConfigPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return homeDir + configFileName, nil
}

func (c *Config) SetUser(username string) error {
	if username == "" {
		return fmt.Errorf("username is an empty string")
	}
	c.CurrentUserName = username
	return write(*c)
}

func write(cfg Config) error {
	configPath, err := getConfigPath()
	if err != nil {
		return err
	}
	jsonFile, err := os.Create(configPath)
	if err != nil {
		return err
	}
	defer jsonFile.Close()
	encoder := json.NewEncoder(jsonFile)
	if err = encoder.Encode(cfg); err != nil {
		return err
	}
	return nil
}
