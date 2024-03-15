package main

import (
	"encoding/json"
	"flag"
	ftpserverlib "github.com/fclairamb/ftpserverlib"
	log "github.com/sirupsen/logrus"
	"hubspot-call_contact/hubspot"
	"strconv"
	"strings"
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

	cfg := hubspot.NewClientConfig()

	// Vernon Communications API Key
	cfg.OAuthToken = config.CrmAPIKey

	c := hubspot.NewClient(cfg)

	/*testCall := hubspot.CallProperties{
		HsTimestamp:       time.Now(),
		HsCallDirection:   "INBOUND",
		HsCallDisposition: "9d9162e7-6cf3-4944-bf63-4dff82258764",
		HsCallDuration:    "666",
		HsCallFromNumber:  "2506665555",
		HsCallStatus:      "BUSY",
		HsCallToNumber:    "123456789",
	}

	err := c.Calls().CreateCall(testCall)
	if err != nil {
		log.Error(err)
		return
	}*/

	owners, err := c.Owners().GetOwners()
	if err != nil {
		log.Error(err)
	}

	err = SaveOwnersToFile(owners.Results, config.CrmUsersFile)
	if err != nil {
		log.Error(err)
	}

	users, err := getZultysUsers(config)
	if err != nil {
		log.Error(err)
	}

	err = SaveUsersToFile(users, config.ZultysUsersFile)
	if err != nil {
		return
	}

	// fmt.Println(string(marshal))

	/*for _, u := range users {
		// todo save the users to a json file
		log.Warn(u)
	}*/

	ch := make(chan DataContent)

	driver := &CustomFtpDriver{
		Username:   config.FtpUsername,
		Password:   config.FtpPassword,
		ListenAddr: config.ListenAddr,
		DataChan:   ch,
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
		t := <-driver.DataChan
		// todo process the lines wether it's from the smdr or ftp upload method
		data, err := processData(t, config)
		if err != nil {
			log.Error(err)
			continue
		}

		ownerFile, err := LoadOwnersFromFile(config.CrmUsersFile)
		if err != nil {
			log.Error(err)
			continue
		}

		for _, d := range data {
			log.Infof("%s", d)
			owner := FindOwnerByEmail(ownerFile, d.User.Email)
			if owner == nil {
				log.Error("could not find owner by email")

				owner = FindOwnerByName(ownerFile, d.User.FirstName, d.User.LastName)
				if owner == nil {
					log.Error("could not find owner by first and last name")
					continue
				}
			}

			log.Infof("found owner %s %s", owner.Email, owner.Id)

			testCall := hubspot.CallProperties{
				HsTimestamp:       d.Time,
				HsCallBody:        "Call imported from Zultys",
				HsCallDirection:   strings.ToUpper(string(d.Direction)),
				HsCallDisposition: "Connected",
				HsCallDuration:    strconv.FormatInt(d.Duration.Milliseconds(), 10),
				HsCallFromNumber:  d.Caller,
				HsCallStatus:      "COMPLETED",
				HsCallToNumber:    d.Callee,
				HsCallTitle:       "Call from " + d.Caller + " to " + d.Callee,
				HubspotOwnerId:    owner.Id,
			}

			marshal, err := json.Marshal(testCall)
			if err != nil {
				return
			}

			println(string(marshal))

			err = c.Calls().CreateCall(testCall)
			if err != nil {
				log.Error(err)
				continue
			}
			// todo create call?? eventually link it to contact as well if one exists...
		}
	}
}
