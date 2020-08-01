package main

import (
	"encoding/json"
	"os"
	"runtime"
)

const (
	configName string = "cfg.json"
)

// Default configurations
var (
	defaultHTTPDisk = &Config{
		MaxWorkers: runtime.GOMAXPROCS(0),
		Crawler:    "http",
		Writer:     "disk",
		Filters: []string{
			"facebook",
			"instagram",
			"google",
			"youtube",
			"amazon",
			"microsoft",
			"apple",
		},
		Paths: []string{
			"https://en.wikipedia.org/wiki/Chaos_Theory",
			"https://en.wikipedia.org/wiki/Machine_Learning",
		},
	}

	defaultFsDisk = &Config{
		MaxWorkers: runtime.GOMAXPROCS(0),
		Crawler:    "fs",
		Writer:     "disk",
		Filters: []string{
			".cache",
			".config",
			".Trash",
		},
		Paths: []string{
			"/home",
		},
	}
)

// Config file
type Config struct {
	MaxWorkers int
	Crawler    string
	Writer     string
	Filters    []string
	Paths      []string
}

// SaveConfig file
func SaveConfig(cf *Config) error {
	f, err := os.Create(configName)
	if err != nil {
		return err
	}
	defer f.Close()

	encoder := json.NewEncoder(f)
	encoder.SetIndent("", "\t")
	encoder.Encode(cf)

	return nil
}

// LoadConfig file
func LoadConfig() (Config, error) {
	var cf Config
	f, err := os.Open(configName)
	if err != nil {
		var cfg *Config

		switch *config {
		case "httpdisk":
			cfg = defaultHTTPDisk
		case "fsdisk":
			cfg = defaultFsDisk
		default:
			cfg = defaultHTTPDisk
		}

		SaveConfig(cfg)
		return cf, err
	}
	defer f.Close()

	decoder := json.NewDecoder(f)
	err = decoder.Decode(&cf)

	return cf, err
}
