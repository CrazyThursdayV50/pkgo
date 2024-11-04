package config

import (
	"errors"
	"log"

	"github.com/spf13/viper"
)

// Parse config file
func parseConfig[Config any](v *viper.Viper) (*Config, error) {
	var c Config
	err := v.Unmarshal(&c)
	if err != nil {
		log.Printf("unable to decode into struct, %v", err)
		return nil, err
	}

	return &c, nil
}

// Load config file from given path
func loadConfig(path, name, ext string) (*viper.Viper, error) {
	v := viper.New()
	v.AddConfigPath(path)
	v.SetConfigName(name)
	v.SetConfigType(ext)
	v.AutomaticEnv()
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil, errors.New("config file not found")
		}
		return nil, err
	}

	return v, nil
}

// Get config
func GetConfig[Config any](path, name, ext string) (*Config, error) {
	cfgFile, err := loadConfig(path, name, ext)
	if err != nil {
		return nil, err
	}

	cfg, err := parseConfig[Config](cfgFile)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}
