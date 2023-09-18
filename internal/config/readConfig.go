package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
)

const (
	pathIsNotSet = "Config path is not set"
	fileNotFound = "Config file not found"
	errorReading = "Error reading config file"
)

func MustReadConfig() Config {
	path := os.Getenv("ConfigPath")
	checkPathIsSet(path)
	checkFileExists(path)
	cfg := readConfig(path)
	return cfg
}

func checkPathIsSet(path string) {
	if path == "" {
		log.Fatal(pathIsNotSet)
	}
}

func checkFileExists(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		log.Println(fileNotFound)
		log.Fatal(err)
	}
}

func readConfig(path string) Config {
	var cfg Config
	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		log.Println(errorReading)
		log.Fatal(err)
	}
	return cfg
}
