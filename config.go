package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml/v2"
)

const (
	APP_VERSION       = "0.0.1"
	CONFIG_FILE_NAME  = "apm.toml"
	BIN_DIR_NAME      = "bin"
	PKG_DIR_NAME      = "packages"
	CLONE_DIR_NAME    = "installers"
	PACKAGE_INDEX_URL = "https://raw.githubusercontent.com/AttifyOS/package-index/main/index.toml"
)

type PackageInfo struct {
	Name        string `toml:"name"`
	Description string `toml:"description"`
	Version     string `toml:"version"`
	RepoUrl     string `toml:"repo_url"`
	RepoTag     string `toml:"repo_tag"`
	License     string `toml:"license"`
}

type APMConfig struct {
	InstallPath string                 `toml:"install_path"`
	Packages    map[string]PackageInfo `toml:"packages"`
}

var _cfg *APMConfig

func (c *APMConfig) LoadFromFile(path string) error {
	cfgBytes, err := os.ReadFile(path)
	if err != nil {
		log.Println("Failed to read config file", path)
		return err
	}
	if err := toml.Unmarshal(cfgBytes, &c); err != nil {
		log.Println("Failed to parse config file", path)
		return err
	}
	return nil
}

func (c *APMConfig) WriteToFile(path string) error {
	cfgMarshalled, err := toml.Marshal(c)
	if err != nil {
		log.Println("Failed to marshal config")
		return err
	}

	if err := os.WriteFile(path, cfgMarshalled, 0644); err != nil {
		log.Println("Failed to write config file", path)
		return err
	}
	return nil
}

func (c *APMConfig) Save() error {
	cfgMarshalled, err := toml.Marshal(c)
	if err != nil {
		log.Println("Failed to marshal config")
		return err
	}
	path := filepath.Join(c.InstallPath, CONFIG_FILE_NAME)
	if err := os.WriteFile(path, cfgMarshalled, 0644); err != nil {
		log.Println("Failed to write config file", path)
		return err
	}
	return nil
}

func GetConfig() *APMConfig {
	if _cfg == nil {
		_cfg = &APMConfig{}
	}
	return _cfg
}
