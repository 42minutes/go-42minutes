package main

import (
	"encoding/json"
	"os"
)

// TraktConfig -
type TraktConfig struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

// RethinkConfig -
type RethinkConfig struct {
	Address string `json:"address"`
}

// Config -
type Config struct {
	Trakt   *TraktConfig   `json:"trakt"`
	Rethink *RethinkConfig `json:"rethinkdb"`
}

func loadConfig(fp string) (*Config, error) {
	fl, err := os.Open(fp)
	if err != nil {
		return nil, err
	}
	defer fl.Close()

	cfg := &Config{}
	decoder := json.NewDecoder(fl)
	err = decoder.Decode(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
