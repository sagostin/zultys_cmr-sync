package main

import (
	"flag"
	log "github.com/sirupsen/logrus"
	"hubspot-call_contact/hubspot"
	"time"
)

func main() {
	log.Infof("%s", "Loading HubSpot Call Associator...")

	// API Key Flags
	var apiKeyFlag = flag.String("apikey", "nil", "specify api key for hubspot")
	if apiKeyFlag == nil {
		log.Fatalf("No HubSpot API key specificed in configuration file.")
	}

	var pRecentCalls = flag.Bool("recent", false, "process only recent calls, "+
		"if no recent processed call exists,"+
		"it will process all calls")

	var storageFile = flag.String("storage", "storage.json", "storage file for latest call processed timestamp")

	flag.Parse()

	apiKey := *apiKeyFlag
	if apiKey == "nil" {
		log.Fatalf("No HubSpot API key specificed in configuration file. 2")
	}

	storage := *storageFile

	// todo recent flag instead of time range
	// recent flag will get the any call newer that the most recently saved call.
	// only when the flag for recent is set will it save the latest call timestamp.

	// Create client object with config from environment variables (HUBSPOT_API_HOST, HUBSPOT_API_KEY, HUBSPOT_OAUTH_TOKEN)
	cfg := hubspot.NewClientConfig()

	// Vernon Communications API Key
	cfg.APIKey = apiKey

	c := hubspot.NewClient(cfg)

	var calls []hubspot.CallResult
	var err error
	bRecentC := *pRecentCalls

	var appStorage AppStorage

	if bRecentC {
		log.Infof("%s", "Running in RECENT call mode.")

		// todo check if config contains a timestamp, else, process all calls and get timestamp
		appStorage, err := LoadAppStorage(storage)
		if err != nil {
			return
		}

		/*err = appStorage.SaveStorage(storage)
		if err != nil {
			log.Fatalf("%s", err)
			return
		}*/

		log.Warnf("%s", time.Now())

		if appStorage.LatestCallTimestamp != (time.Time{}) {
			log.Infof("%s %s", "Found recent call association timestamp... Processing only recent calls.", appStorage.LatestCallTimestamp)
			recentDuration := time.Now().Sub(appStorage.LatestCallTimestamp)
			calls, err = c.Calls().GetRecentCalls(recentDuration)
			if err != nil {
				log.Fatalf("%s", err)
			}
		} else {
			log.Infof("%s", "No latest call found... Processing all calls...")
			calls, err = c.Calls().GetAllCalls()
			if err != nil {
				log.Fatalf("%s", err)
			}
		}
	} else {
		log.Infof("%s", "Processing all calls...")
		calls, err = c.Calls().GetAllCalls()
		if err != nil {
			log.Fatalf("%s", err)
		}
	}

	log.Warnf("Sleeping for 5 seconds to prevent rate limits...\n" +
		"------------------------------------------------------------")
	time.Sleep(5 * time.Second)

	var latestCall time.Time

	for _, r := range calls {
		callNum := ""
		if r.Properties.HsCallDirection == "INBOUND" {
			callNum = r.Properties.HsCallFromNumber
		} else if r.Properties.HsCallDirection == "OUTBOUND" {
			callNum = r.Properties.HsCallToNumber
		}

		if callNum == "" {
			log.Warnf("Skipping... No number found...")
			sleepRateLimit()
			continue
		}

		contacts, err := c.Contacts().SearchByPhone(callNum)
		if err != nil {
			log.Error(err)
			sleepRateLimit()
			continue
		}

		if len(contacts.Results) == 0 {
			log.Warnf("No contacts found... Skipping...")
			sleepRateLimit()
			continue
		}

		foundContact := false
		for _, cr := range contacts.Results {
			if r.Properties.HsCallCalleeObjectId == "" {
				err := c.Calls().AssociateCallContact(r, cr, 194)
				if err != nil {
					log.Warnf("%s", err)
				}
			}

			log.Infof("*CONTACT* Name: %s %s, Company: %s Phone: %s Email: %s", cr.Properties.FirstName,
				cr.Properties.LastName, cr.Properties.Company,
				cr.Properties.Phone, cr.Properties.Email)

			foundContact = true

			// check if call contains hs_call_callee_object_id
			// if it doesn't and theres contacts available, associate the calls.
		}

		if foundContact {
			if &latestCall == nil {
				latestCall = r.Properties.HsTimestamp
				log.Warnf("Latest Call with Contact: %s", latestCall.String())
			} else if r.Properties.HsTimestamp.After(latestCall) {
				latestCall = r.Properties.HsTimestamp
				log.Warnf("Latest Call with Contact: %s", latestCall.String())
			}
		}

		sleepRateLimit()
	}
	if latestCall != (time.Time{}) {
		appStorage.LatestCallTimestamp = latestCall
		err = appStorage.SaveStorage("storage.json")
		if err != nil {
			log.Fatalf("Unable to save storage.json : %s", err)
			return
		}
	}
	log.Infof("Finished running HubSpot Call Associator")
}

func sleepRateLimit() {
	// Wait for 10 second time out just in case
	//log.Warnf("Sleeping for 1 second to prevent rate limits...\n")
	time.Sleep(1 * time.Second)
}
