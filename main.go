package main

import (
	"context"
	"errors"
	"flag"
	"go.guoyk.net/sugar/sugar_zap"
	"go.guoyk.net/unizone/pkg/providers"
	"os"
	"path/filepath"
	"strings"
)

const PrefixCached = "CACHED "

func main() {
	var (
		optVerbose bool
		optConf    string
		optOutput  string
	)

	flag.BoolVar(&optVerbose, "v", false, "enable verbose logging")
	flag.StringVar(&optConf, "c", "", "configuration file")
	flag.StringVar(&optOutput, "o", "", "output zone file")
	flag.Parse()

	log := createLogger(optVerbose)

	var err error
	defer func(err *error) {
		_ = log.Sync()
		if *err != nil {
			os.Exit(1)
		}
	}(&err)

	if optConf == "" {
		err = errors.New("missing argument -c")
		log.Error(err.Error())
		return
	}

	if optOutput == "" {
		err = errors.New("missing argument -o")
		log.Error(err.Error())
		return
	}

	var cfg Config
	if cfg, err = loadConfig(optConf); err != nil {
		log.Errorf("failed to load config file: %s, %s", optConf, err.Error())
		return
	}

	cacheDir := filepath.Join(filepath.Dir(optConf), cfg.CacheDir)

	var services []providers.Service

	for _, cloud := range cfg.Clouds {
		for _, region := range cloud.Regions {
			var provider providers.Provider
			if provider, err = providers.Create(cloud.Provider, providers.Options{
				TokenID:     cloud.TokenID,
				TokenSecret: cloud.TokenSecret,
				Region:      region,
				Logger:      sugar_zap.Wrap(log.With("cloud", cloud.Name, "region", region).Desugar()),
			}); err != nil {
				log.Warnf("failed to create provider: %s (%s/%s), %s", cloud.Name, cloud.Provider, region, err.Error())
				err = nil
				continue
			}
			for _, network := range cloud.Networks {
				for _, service := range cloud.Services {
					var cloudServices []providers.Service
					if cloudServices, err = provider.ListServices(context.Background(), network, service); err != nil {
						log.Warnf("failed to list services: %s (%s/%s) %s/%s, %s", cloud.Name, cloud.Provider, region, network, service, err.Error())
						err = nil
						continue
					}
				outerLoop1:
					for _, cloudService := range cloudServices {
						for _, knownService := range services {
							if knownService.Name == cloudService.Name {
								log.Warnf("found duplicated service name: %s (%s/%s) %s/%s/%s, %s (%s) vs %s", cloud.Name, cloud.Provider, region, network, service, cloudService.Name, cloudService.IP, knownService.Comment, err.Error())
								continue outerLoop1
							}
						}
						log.Debugf("found service: %s (%s) %s/%s/%s (%s)", cloud.Name, cloud.Provider, network, services, cloudService.Name, cloudService.IP)
						services = append(services, cloudService)
					}
				}
			}
		}
	}

	var cachedServices []providers.Service
	if cachedServices, err = loadCachedServices(cacheDir); err != nil {
		log.Errorf("failed to load cached services from dir: %s, %s", cacheDir, err.Error())
		return
	}

outerLoop2:
	for _, cachedService := range cachedServices {
		for _, knownService := range services {
			if knownService.Name == cachedService.Name {
				continue outerLoop2
			}
		}
		log.Warnf("cached service appended: %s %s (%s)", cachedService.Comment, cachedService.Name, cachedService.IP)
		if !strings.HasPrefix(cachedService.Comment, PrefixCached) {
			cachedService.Comment = PrefixCached + cachedService.Comment
		}
		services = append(services, cachedService)
	}

	if err = saveCachedServices(services, cacheDir); err != nil {
		log.Errorf("failed to save services to cache dir: %s, %s", cacheDir, err.Error())
		return
	}

	if err = writeZoneFile(services, optOutput); err != nil {
		return
	}
}
