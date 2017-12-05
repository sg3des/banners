package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

func findConfigFile() string {
	files := []string{
		filepath.Join("testdata", "banners.conf"),
		"/etc/banners.conf",
		"banners.conf",
	}

	for _, filename := range files {
		fi, err := os.Stat(filename)
		if err == nil && !fi.IsDir() {
			return filename
		}
	}

	return ""
}

//LoadConfig function read config file and decode it, config is should be in toml format
func LoadConfig(filename string) error {
	if filename == "" {
		filename = findConfigFile()
	}

	if filename == "" {
		return errors.New("configuration file not found")
	}

	setDefaultValues()

	_, err := toml.DecodeFile(filename, &args)
	if err != nil {
		return fmt.Errorf("failed load configuration from file '%s', reason: %s", filename, err)
	}

	return nil
}

func setDefaultValues() {
	if args.HTTPAddr == "" {
		args.HTTPAddr = ":80"
	}
	if args.CSVFile == "" {
		args.CSVFile = "banners.csv"
	}
	if args.CSVSeparator == "" {
		args.CSVSeparator = ";"
	}
	if args.StorageType == "" {
		args.StorageType = "chan"
	}
}
