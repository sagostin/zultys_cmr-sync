package main

import (
	"encoding/json"
	ftpserverlib "github.com/fclairamb/ftpserverlib"
	log "github.com/sirupsen/logrus"
	"hubspot-call_contact/hubspot"
	"os"
	"regexp"
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
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "./config.json"
	}
	config := loadConfig(configPath)

	cfg := hubspot.NewClientConfig()

	// Vernon Communications API Key
	cfg.OAuthToken = config.CrmAPIKey

	c := hubspot.NewClient(cfg)

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

		go func() {
			ownerFile, err := LoadOwnersFromFile(config.CrmUsersFile)
			if err != nil {
				log.Error(err)
				return
			}

			for _, d := range data {
				log.Infof("%s", d)
				owner := FindOwnerByEmail(ownerFile, d.User.Email)
				if owner == nil {
					log.Error("could not find owner by email")

					owner = FindOwnerByName(ownerFile, d.User.FirstName, d.User.LastName)
					if owner == nil {
						log.Error("could not find owner by first and last name")

						// create empty owner
						owner = &hubspot.Owner{}
						//continue
					}
				}

				log.Infof("found owner %s %s", owner.Email, owner.Id)

				// todo associate the call with companies & contacts??

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

				callResult, err := c.Calls().CreateCall(testCall)
				if err != nil {
					log.Error(err)
					continue
				}

				/*_, err = c.Calls().GetAllCalls()
				if err != nil {
					return
				}*/

				/*call, err := c.Calls().GetCall(callResult.Id)
				if err != nil {
					log.Error(err)
				}*/

				var phoneLookupNum string
				if d.Direction == OUTBOUND {
					phoneLookupNum = d.Callee
				} else if d.Direction == INBOUND {
					phoneLookupNum = d.Caller
				}

				match10Digit, _ := regexp.MatchString("^(\\d{10})$", phoneLookupNum)
				match11Digit, _ := regexp.MatchString("^(1\\d{10})$", phoneLookupNum)
				if match10Digit {
					phoneLookupNum = "+1" + phoneLookupNum
				} else if match11Digit {
					phoneLookupNum = "+" + phoneLookupNum
				}

				companies, err := c.Companies().SearchByPhone(phoneLookupNum)
				if err != nil {
					log.Error(err)
				}

				for _, companyP := range companies.Results {
					err = c.Calls().AssociateCallCompany(*callResult, companyP.Id, 182)
					if err != nil {
						log.Warn(err)
					}
				}

				phoneContact, err := c.Contacts().SearchByPhone(phoneLookupNum)
				if err != nil {
					log.Error(err)
				}

				for _, contactP := range phoneContact.Results {
					err = c.Calls().AssociateCallContact(*callResult, contactP, 194)
					if err != nil {
						log.Warn(err)
					}
					// todo create call?? eventually link it to contact as well if one exists...
				}
			}
		}()
	}
}
