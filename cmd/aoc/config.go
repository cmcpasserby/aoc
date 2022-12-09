package main

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"io"
	"os"
	"path/filepath"
	"time"
)

const (
	configName = "aocConfig.toml"
)

type aocConfig struct {
	SessionCookie string `toml:"session_cookie"`
	Year          int    `toml:"year"`
	Day           int    `toml:"day"`
}

func getConfig() (*aocConfig, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	workingDirPath := filepath.Join(cwd, configName)

	f, err := os.Open(workingDirPath)
	if err == nil {
		defer f.Close()
		return loadConfig(f)
	}

	configDir, err := os.UserConfigDir()
	if err != nil {
		return nil, err
	}
	userConfigPath := filepath.Join(configDir, configName)

	f, err = os.Open(userConfigPath)
	if err == nil {
		defer f.Close()
		return loadConfig(f)
	}

	return nil, fmt.Errorf("no aocConfig found in '%s' or '%s'", workingDirPath, userConfigPath)
}

func loadConfig(f io.Reader) (*aocConfig, error) {
	var config aocConfig
	if _, err := toml.DecodeReader(f, &config); err != nil {
		return nil, err
	}
	return &config, nil
}

func getCombinedConfig() (*aocConfig, error) {
	config, err := getConfig()
	if err != nil {
		return nil, err
	}

	if gFlagSessionCookie != "" {
		config.SessionCookie = gFlagSessionCookie
	}

	loc, err := time.LoadLocation("America/New_York")
	if err != nil {
		return nil, err
	}

	now := time.Now().In(loc)
	year := now.Year()
	if gFlagYear != 0 {
		year = gFlagYear
	} else if config.Year != 0 {
		year = config.Year
	}
	config.Year = year

	day := now.Day()
	if gFlagDay != 0 {
		day = gFlagDay
	} else if config.Day != 0 {
		day = config.Day
	}
	config.Day = day

	return config, nil
}
