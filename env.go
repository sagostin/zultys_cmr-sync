package main

import (
	"errors"
	log "github.com/sirupsen/logrus"
	"os"
)

const (
	defaultConfig = "{\n \"latest_call_timestamp\": \"2000-01-01T00:00:00.0000Z\"\n}"
)

func setup(storageFile string) error {
	// Check if the config file exists in the local directory
	_, err := os.Stat(storageFile)
	// If the check returns an error indicating the file doesn't exist, create it
	if errors.Is(err, os.ErrNotExist) {
		// Log to terminal that a new file will be created
		log.Warnf("Config file does '%s' does not exist, creating one now.", storageFile)
		// Attempt to create the config file
		_, err = os.Create(storageFile)
		if err != nil {
			return err
		}
		// Attempt to write the default config pattern to the config file
		err = os.WriteFile("./storage.json", []byte(defaultConfig), 0644)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}
	return nil
}
