package main

import (
	"flag"
	ftpserverlib "github.com/fclairamb/ftpserverlib"
	log "github.com/sirupsen/logrus"
)

/*

- Start with running local FTP server and/or SMDR receiver server.
- Provide json configuration for the following:
	- MX Hostname
	- MX Admin Username
	- MX Admin Password
	* We need the MX login information to be able to pull the extension list and correlate the users/extensions
	? Should we resolve the MX Hostname and block all further attempts to the FTP server from other IPs from within
	the application itself
	- Sync Mode (SMDR or MXReport FTP based)
	? FTP will receive files in memory? then process the incoming files and then go line by line and match the specified
	regex format
	- FTP Host & Port
	- SMDR Host & Port
	? Can we just have a general listen field and use that for both? We will run this application as a dedicated instance
	for both types
	- CMR (this will be hubspot to begin with, Zoho in future)
*/

func main() {
	configPath := flag.String("config", "./config.json", "config file path")
	flag.Parse()
	config := loadConfig(*configPath)

	ch := make(chan string)

	driver := &CustomFtpDriver{
		Username:    config.FtpUsername,
		Password:    config.FtpPassword,
		ListenAddr:  config.ListenAddr,
		CsvDataChan: ch,
	}

	// Instantiate the FTP server using our custom driver
	go func() {
		log.Infof("loading ftp server thread")
		server := ftpserverlib.NewFtpServer(driver)

		// Start the server
		log.Infof("Starting FTP server on %s...", driver.ListenAddr)
		if err := server.ListenAndServe(); err != nil {
			log.Fatal("Error starting server: ", err)
		}
	}()

	for {
		t := <-driver.CsvDataChan
		// todo process the lines wether it's from the smdr or ftp upload method
		log.Warn(t)
	}
}
