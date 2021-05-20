package main

import (
	"go.guoyk.net/unizone/pkg/providers"
	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

func createLogger(verbose bool) *zap.SugaredLogger {
	cfg := zap.NewDevelopmentConfig()
	if verbose {
		cfg.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
		cfg.Development = true
	} else {
		cfg.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
		cfg.Development = false
	}
	cfg.Encoding = "console"
	logger, _ := cfg.Build()
	return logger.Sugar()
}

func loadConfig(file string) (cfg Config, err error) {
	var buf []byte
	if buf, err = ioutil.ReadFile(file); err != nil {
		return
	}
	if err = yaml.Unmarshal(buf, &cfg); err != nil {
		return
	}
	return
}

func loadCachedServices(dir string) (services []providers.Service, err error) {
	// TODO:
	return
}

func saveCachedServices(services []providers.Service, dir string) (err error) {
	// TODO:
	return
}

func writeZoneFile(services []providers.Service, file string) (err error) {
	// TODO:
	return
}
