package main

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Config struct {
	Origin    string `yaml:"origin"`
	TTL       int    `yaml:"ttl"`
	Providers []struct {
		ID          string `yaml:"id"`
		Provider    string `yaml:"provider"`
		TokenID     string `yaml:"token_id"`
		TokenSecret string `yaml:"token_secret"`
		Networks    []struct {
			ID     string `yaml:"id"`
			Region string `yaml:"region"`
		} `yaml:"networks"`
		Services []string `yaml:"services"`
	} `yaml:"providers"`
}

func LoadConfigFile(file string, cfg *Config) (err error) {
	var buf []byte
	if buf, err = ioutil.ReadFile(file); err != nil {
		return
	}
	if err = yaml.Unmarshal(buf, cfg); err != nil {
		return
	}
	if cfg.Origin == "" {
		err = errors.New("missing 'origin' field")
		return
	}
	if cfg.TTL < 0 {
		err = errors.New("invalid 'ttl' field")
		return
	}
	if len(cfg.Providers) == 0 {
		err = errors.New("missing 'providers' field")
		return
	}
	for i, pvd := range cfg.Providers {
		if pvd.ID == "" {
			err = fmt.Errorf("missing 'id' field in providers.%d", i+1)
			return
		}
		if pvd.Provider == "" {
			err = fmt.Errorf("missing 'provider' field in providers.%d", i+1)
			return
		}
		if len(pvd.Networks) == 0 {
			err = fmt.Errorf("missing 'networks' field in providers.%d", i+1)
			return
		}
		for j, network := range pvd.Networks {
			if network.Region == "" {
				err = fmt.Errorf("missing 'region' field in providers.%d.networks.%d", i+1, j+1)
				return
			}
			if network.ID == "" {
				err = fmt.Errorf("missing 'id' field in providers.%d.networks.%d", i+1, j+1)
				return
			}
		}
		if len(pvd.Services) == 0 {
			err = fmt.Errorf("missing 'services' field in providers.%d", i+1)
			return
		}
	}
	return
}
