package config

// from: https://dev.to/ilyakaznacheev/a-clean-way-to-pass-configs-in-a-go-application-1g64

import (
	"errors"
	"fmt"
	"os"

	"github.com/kelseyhightower/envconfig"
	"gopkg.in/yaml.v2"
)

var cfg Config

func GetConfig() *Config {
	return &cfg
}

func processError(err error) {
	fmt.Println(err)
}

func first[T, U any](val T, _ U) T {
	return val
}

func exists(name string) (bool, error) {	
	_, err := os.Stat(name)
	if err == nil {
		return true, nil
	}
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	return false, err
}

func GetConfigPath(filepath string) (string, bool) {
	var config string

	if first(exists(filepath)) {
		return filepath, true
	}

	config = os.Getenv("FUSEFS-CONFIG")
	if config != "" {
		return config, true
	}

	return "", false
}

func readFile(cfg *Config, filename string) bool {
	cfgPath, isExist := GetConfigPath(filename)

	if !isExist {
		return false
	}

	f, err := os.Open(cfgPath)
	if err != nil {
		processError(err)
	}
	defer f.Close()

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(cfg)
	if err != nil {
		processError(err)
	}

	return true
}

func readEnv(cfg *Config) {
	err := envconfig.Process("", cfg)

	if err != nil {
		processError(err)
	}
}

func BuildConfig(filename string) {

	if !readFile(&cfg, filename) {
		fmt.Println("config: config file os not found, using just environment variables")
	}
	readEnv(&cfg)
}
