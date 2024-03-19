package main

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"os"
)

type Config struct {
	MxUsername             string   `json:"mx_username"`
	MxPassword             string   `json:"mx_password"`
	MxAddr                 string   `json:"mx_addr"`
	ListenAddr             string   `json:"listen_addr"`
	Mode                   DataType `json:"mode"`
	FtpUsername            string   `json:"ftp_username"`
	FtpPassword            string   `json:"ftp_password"`
	CrmType                string   `json:"crm_type"`
	CrmAPIKey              string   `json:"crm_apikey"`
	ZultysUsersFile        string   `json:"zultys_users_file"`
	CrmUsersFile           string   `json:"crm_users_file"`
	TimestampRegion        string   `json:"timestamp_region"`
	TimestampFile          string   `json:"timestamp_file"`
	IncludeUnknownCRMUsers bool     `json:"include_unknown_crm_users"`
}

func loadConfig(configFilePath string) Config {
	var config Config
	// Open the configuration file
	file, err := os.Open(configFilePath)
	if err != nil {
		log.Errorf("Error opening config file: %v\n", err)
		return config
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Error(err)
		}
	}(file)

	// Decode the configuration file into the Config struct
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		log.Errorf("Error decoding config file: %v\n", err)
		return config
	}

	// Now you can use your config struct
	log.Infof("Loaded configuration: %+v\n", config)
	return config
}
