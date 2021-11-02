package main

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"io"
	"os"
	"path/filepath"
)

const (
	configName = ".aocConfig"
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

	configDir, err := os.UserHomeDir()
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
