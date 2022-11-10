package main

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"os"
	"time"
)

type AppStorage struct {
	//ApiKey              string    `json:"api_key"`
	LatestCallTimestamp time.Time `json:"latest_call_timestamp"`
}

func (c AppStorage) SaveStorage(configPath string) error {
	file, err := os.OpenFile(configPath, os.O_CREATE|os.O_APPEND|os.O_RDWR|os.O_TRUNC, 0644)

	err = file.Truncate(0)
	_, err = file.Seek(0, 0)

	if err != nil {
		log.Errorf("")
		return err
	}

	f, _ := json.MarshalIndent(c, "", " ")
	file.Write(f)

	return file.Close()
}

func LoadAppStorage(configPath string) (*AppStorage, error) {
	file, err := os.ReadFile(configPath)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	//log.Infof("")

	var c *AppStorage

	err = json.Unmarshal(file, &c)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	return c, nil
}
