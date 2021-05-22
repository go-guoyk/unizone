package main

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Config struct {
	CacheDir string `yaml:"cache_dir"`
	Origin   string `yaml:"origin"`
	TTL      int    `yaml:"ttl"`
	Clouds   []struct {
		Name        string   `yaml:"name"`
		Provider    string   `yaml:"provider"`
		TokenID     string   `yaml:"token_id"`
		TokenSecret string   `yaml:"token_secret"`
		Regions     []string `yaml:"regions"`
		Networks    []string `yaml:"networks"`
		Services    []string `yaml:"services"`
	} `yaml:"clouds"`
}

func LoadConfigFile(file string, cfg *Config) (err error) {
	var buf []byte
	if buf, err = ioutil.ReadFile(file); err != nil {
		return
	}
	if err = yaml.Unmarshal(buf, cfg); err != nil {
		return
	}
	return
}
