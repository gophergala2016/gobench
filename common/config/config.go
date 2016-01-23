package config

import (
	"encoding/json"
	"github.com/gophergala2016/gobench/backend"
	"github.com/gophergala2016/gobench/frontend"
	"io/ioutil"
)

type Config struct {
	Backend  backend.Config  `json:"backend"`
	Frontend frontend.Config `json:"frontend"`
}

// ReadFromFile reads application configuration from JSON file
func ReadFromFile(fname string) (*Config, error) {
	c := Config{}
	buf, err := ioutil.ReadFile(fname)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(buf, &c); err != nil {
		return nil, err
	}

	return &c, nil
}

// ReadFromEnvVars reads application configuration fron environment variables
func ReadFromEnvVars(prefix string) (*Config, error) {
	// TODO. Probable after GopherGala
	c := Config{}
	return &c, nil
}
