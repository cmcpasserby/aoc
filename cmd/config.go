package cmd

import (
	"github.com/BurntSushi/toml"
	"os"
	"path/filepath"
)

const configName = "config.toml"

type config struct {
	SessionId string `toml:"session_id"`
}

func loadConfig() (*config, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	path := filepath.Join(cwd, configName)

	var config config
	_, err = toml.DecodeFile(path, &config)
	return &config, err
}
