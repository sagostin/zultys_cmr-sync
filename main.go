package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	ftpserverlib "github.com/fclairamb/ftpserverlib"
	"github.com/sagostin/zultys_crm-sync/hubspot"
	"github.com/sagostin/zultys_crm-sync/zoho"
	log "github.com/sirupsen/logrus"
	"net"
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
	// Define the command-line flag
	configPathFlag := flag.String("config", "./config.json", "path to config file")
	flag.Parse()

	// Define the environment variable key
	envVar := "CONFIG_PATH"

	// Check if the environment variable is set
	configPathEnv, envSet := os.LookupEnv(envVar)

	// Use the environment variable if it's set, otherwise use the flag value
	configPath := *configPathFlag
	if envSet {
		configPath = configPathEnv
	}

	// Use configPath as needed
	log.Infof("Using configuration file at: %s\n", configPath)
	config := loadConfig(configPath)

	users, err := getZultysUsers(config)
	if err != nil {
		log.Error(err)
	}

	err = SaveUsersToFile(users, config.ZultysUsersFile)
	if err != nil {
		return
	}

	var c hubspot.Client
	var zohoClient zoho.Client

	if config.CrmType == ZohoCRM {
		zohoStrings := strings.Split(config.CrmAPIKey, ":")

		zohoClient = zoho.Client{Endpoints: zoho.Endpoints{
			AccountAuth: zohoStrings[0],
			CrmApi:      zohoStrings[1],
		},
			Auth: zoho.AccessGrant{},
		}

		err := zohoClient.Authenticate(zohoStrings[2], zohoStrings[3], zohoStrings[4])
		if err != nil {
			log.Error(err)
			return
		}

		zohoClient.StartTokenRefresher()

		log.Info("authenticated with zoho")

		fetchUsers, err := zohoClient.FetchUsers()
		if err != nil {
			log.Error(err)
			return
		}

		err = zoho.SaveUsersToFile(fetchUsers.Users, config.CrmUsersFile)
	} else if config.CrmType == HubspotCRM {

		cfg := hubspot.NewClientConfig()

		// Vernon Communications API Key
		cfg.OAuthToken = config.CrmAPIKey

		c = hubspot.NewClient(cfg)

		// todo periodic refresh/update of the file to always be up to date???

		owners, err := c.Owners().GetOwners()
		if err != nil {
			log.Error(err)
		}

		err = hubspot.SaveOwnersToFile(owners.Results, config.CrmUsersFile)
		if err != nil {
			log.Error(err)
		}
	}

	// fmt.Println(string(marshal))

	/*for _, u := range users {
		// todo save the users to a json file
		log.Warn(u)
	}*/

	ch := make(chan DataContent)

	if config.Mode == DataTypeFTP {
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
	} else if config.Mode == DataTypeSMDR {
		go func(c chan DataContent) {
			// Listen on the specified port
			listener, err := net.Listen("tcp", config.ListenAddr)
			if err != nil {
				fmt.Println("Error listening:", err.Error())
				os.Exit(1)
			}
			defer func(listener net.Listener) {
				err := listener.Close()
				if err != nil {
					log.Error(err)
				}
			}(listener)
			fmt.Println("Listening for SMDR on " + config.ListenAddr)

			// Accept connections in a loop
			for {
				conn, err := listener.Accept()
				if err != nil {
					fmt.Println("Error accepting: ", err.Error())
					os.Exit(1)
				}
				fmt.Println("Connection accepted.")

				// Handle connections in a new goroutine.
				go handleSmdrRequest(conn, c)
			}
		}(ch)
	}

	for {
		content := <-ch
		// todo process the lines wether it's from the smdr or ftp upload method
		go func(t DataContent, zc *zoho.Client, hsc *hubspot.Client) {
			data, err := processData(t, config)
			if err != nil {
				log.Error(err)
				return
			}

			if config.CrmType == ZohoCRM {
				handleZoho(zc, config, data)
			} else if config.CrmType == HubspotCRM {
				handleHubspot(&c, config, data)
			}
		}(content, &zohoClient, &c)
	}
}

func handleZoho(c *zoho.Client, config Config, data []CallEntry) {
	// todo
	ownerFile, err := zoho.LoadUsersFromFile(config.CrmUsersFile)
	if err != nil {
		log.Error(err)
		return
	}

	for _, d := range data {
		log.Infof("%s", d)
		owner := zoho.FindUserByEmail(ownerFile, d.User.Email)
		if owner == nil {
			log.Error("could not find owner by email")

			owner = zoho.FindUserByName(ownerFile, d.User.FirstName, d.User.LastName)
			if owner == nil {
				log.Error("could not find owner by first and last name")

				if config.IncludeUnknownCRMUsers {
					// create empty owner
					owner = &zoho.User{}
				} else {
					continue
				}
			}
		}

		log.Infof("found owner %s %s", owner.Email, owner.Id)

		callStatus := "Completed"

		var phoneLookupNum string
		if d.Direction == OUTBOUND {
			phoneLookupNum = d.Callee
		} else if d.Direction == INBOUND {
			phoneLookupNum = d.Caller
			callStatus = ""
		}

		/*match10Digit, _ := regexp.MatchString("^(\\d{10})$", phoneLookupNum)
		match11Digit, _ := regexp.MatchString("^(1\\d{10})$", phoneLookupNum)
		if match10Digit {
			phoneLookupNum = "+1" + phoneLookupNum
		} else if match11Digit {
			phoneLookupNum = "+" + phoneLookupNum
		}*/

		var WhoId *zoho.ItemLink

		contact, err := c.FindContactByPhone(phoneLookupNum)
		if err != nil {
			log.Errorf("error finding contact by phone: %s", err.Error())
		}

		// todo account for multiple contacts?
		if len(contact.Data) > 0 {
			WhoId.Id = contact.Data[0].Id
		} else {
			log.Error("failed to find contact by phone")
		}

		var WhatId *zoho.ItemLink

		account, err := c.FindAccountByPhone(phoneLookupNum)
		if err != nil {
			log.Errorf("error finding account by phone: %s", err.Error())
		}

		// todo account for multiple accounts?
		if len(account.Data) > 0 {
			WhatId.Id = account.Data[0].Id
		} else {
			log.Error("failed to find account by phone")
		}

		// Convert the duration to total seconds for ease of calculation
		totalSeconds := int(d.Duration.Seconds())

		// Calculate minutes and seconds
		minutes := totalSeconds / 60
		seconds := totalSeconds % 60

		// Format the string into mm:ss format
		durationStr := fmt.Sprintf("%02d:%02d", minutes, seconds)

		var ownerID *zoho.Owner

		if owner != nil && owner.Id != "" {
			ownerID = &zoho.Owner{Id: owner.Id}
		}

		var call zoho.Call

		call = zoho.Call{
			Owner:                 ownerID,
			CallDuration:          durationStr,
			Description:           "Imported from Zultys",
			CallPurpose:           "-None-",
			OutgoingCallStatus:    callStatus,
			CallStartTime:         d.Time,
			Subject:               "Call from " + d.Caller + " to " + d.Callee,
			SeModule:              "Accounts",
			CallType:              string(d.Direction),
			WhatId:                WhatId,
			WhoId:                 WhoId,
			CallDurationInSeconds: d.Duration.Seconds(),
			CreatedBy:             ownerID,
			ModifiedBy:            ownerID,
		}

		marshal, err := json.Marshal(call)
		if err != nil {
			log.Error(err)
		}

		println(string(marshal))

		var calls []zoho.Call
		calls = append(calls, call)

		err = c.CreateCall(calls)
		if err != nil {
			log.Error(err)
			continue
		}
	}
}

func handleHubspot(c *hubspot.Client, config Config, data []CallEntry) {
	// todo
	ownerFile, err := hubspot.LoadOwnersFromFile(config.CrmUsersFile)
	if err != nil {
		log.Error(err)
		return
	}

	for _, d := range data {
		log.Infof("%s", d)
		owner := hubspot.FindOwnerByEmail(ownerFile, d.User.Email)
		if owner == nil {
			log.Error("could not find owner by email")

			owner = hubspot.FindOwnerByName(ownerFile, d.User.FirstName, d.User.LastName)
			if owner == nil {
				log.Error("could not find owner by first and last name")

				if config.IncludeUnknownCRMUsers {
					// create empty owner
					owner = &hubspot.Owner{}
				} else {
					continue
				}
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
}

// Handles incoming requests
func handleSmdrRequest(conn net.Conn, ch chan DataContent) {
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			log.Error(err)
		}
	}(conn)

	// Create a new reader, assuming carriage returns and line feeds as delimiters
	reader := bufio.NewReader(conn)

	// Read data line by line
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading:", err.Error())
			break
		}
		// Trim the line to remove the delimiter
		line = strings.Trim(line, "\r\n")

		// Split the line by spaces to get individual fields
		fields := strings.Split(line, " ")

		// Process the fields (for now, just print them out)
		fmt.Println(fields)

		ch <- DataContent{
			FilePath: "SMDR not FILE",
			Type:     DataTypeSMDR,
			Content:  line,
		}

		// Send a response back to the client (optional)
		_, err = conn.Write([]byte("Received Line\n"))
		if err != nil {
			log.Error(err)
		}
	}
}
