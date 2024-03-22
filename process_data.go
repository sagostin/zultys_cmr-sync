package main

import (
	"encoding/json"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"strings"
	"time"
)

type DataType string

const DataTypeFTP DataType = "ftp"
const DataTypeSMDR DataType = "smdr"

type DataContent struct {
	FilePath string
	Type     DataType
	Content  string
}

type CallDirection string

const INBOUND CallDirection = "inbound"
const OUTBOUND CallDirection = "outbound"
const INTERNAL CallDirection = "internal"

type CallEntry struct {
	Time      time.Time     `json:"time"`
	Direction CallDirection `json:"direction,omitempty"`
	User      ZultysUser
	Extension string
	Duration  time.Duration
	DialPlan  string
	Caller    string
	Callee    string
}

func processData(content DataContent, config Config) ([]CallEntry, error) {
	if content.Type == DataTypeFTP {
		data, err := processFtpData(content, config)
		if err != nil {
			return data, err
		}
		return data, nil
	} else if content.Type == DataTypeSMDR {
		// todo handle smdr data format?
		data, err := processSmdrData(content, config)
		if err != nil {
			return data, err
		}
		return data, nil
	}
	return nil, nil
}

func processSmdrData(content DataContent, config Config) ([]CallEntry, error) {
	var calls []CallEntry

	if strings.Contains(content.FilePath, "SMDR not FILE") {

		file, err := LoadUsersFromFile(config.ZultysUsersFile)
		if err != nil {
			log.Errorf("could not find user? %s", err)
		}

		// Load previously processed timestamps
		extensionTimestamps, err := LoadTimestampsFromFile(config.TimestampFile)
		if err != nil {
			log.Errorf("unable to get timestamps from file, assuming new: %v", err)
			extensionTimestamps = make(ExtensionTimestamps)
		}

		items := strings.Split(content.Content, " ")

		// outbound:
		// N 001 00 DN1092 TVoiceGroup1001 12/16 09:00:56 00:00:03 4155551212
		//
		// inbound:
		// N 002 00 TVoiceGroup1002 DN1093 12/16 09:01:45 00:00:04 4155551212

		if len(items) < 9 || len(items) > 9 {
			return nil, errors.New("invalid data format")
		}

		/*state := items[0]
		recordseq := items[1]
		custnum := items[2]*/
		extension_or_trunk := items[3]
		trunk_or_extension := items[4]
		date := items[5]
		date_time := items[6]
		call_duration := items[7]
		external_party := items[8]

		/*dialplan := values[6]
		if strings.Contains(dialplan, "park") {
			log.Info("skipping over park entry")
			continue
		}*/

		var extension string
		var trunk string

		var direction = INTERNAL
		if strings.HasPrefix(extension_or_trunk, "DN") {
			direction = OUTBOUND
			extension = strings.ReplaceAll(extension_or_trunk, "DN", "")
			trunk = trunk_or_extension
		} else if strings.HasPrefix(extension_or_trunk, "T") {
			direction = INBOUND
			extension = strings.ReplaceAll(trunk_or_extension, "DN", "")
			trunk = extension_or_trunk
		}

		if direction == INTERNAL {
			return nil, nil
		}

		const dateTimeLayout = "1/2/2006 15:04:05" // Combined layout for parsing both together

		// Combine the date and time strings
		dateTimeStr := fmt.Sprintf("%s/%d %s", date, time.Now().Year(), date_time)

		location, err := time.LoadLocation(config.TimestampRegion)
		if err != nil {
			log.Fatalf("Error loading location: %v\n", err)
		}

		// Parse the combined date and time string into a time.Time object
		callTime, err := time.ParseInLocation(dateTimeLayout, dateTimeStr, location)
		if err != nil {
			log.Errorf("Error parsing date and time: %v\n", err)
		}

		// Check if the call is newer than the last processed for this extension
		lastProcessed, exists := extensionTimestamps[extension]
		if exists && callTime.Unix() <= lastProcessed {
			// This call has already been processed, skip it
			log.Warn("skipping call, due to being already processed?")
			return nil, nil
		}

		durationSep := strings.Split(call_duration, ":")
		duration, err := time.ParseDuration(durationSep[0] + "h" + durationSep[1] + "m" + durationSep[2] + "s")
		if err != nil {
			return nil, err
		}

		// fmt.Println(line)

		caller := extension
		callee := external_party

		if direction == INBOUND {
			caller = external_party
			callee = extension
		}

		user := FindEntryByExtension(file, extension)
		if user == nil {
			log.Error("no user found")
			user = &ZultysUser{}
			return nil, errors.New("unable to find matching user in system, skipping")
		}

		entry := CallEntry{
			Time:      callTime,
			Direction: direction,
			User:      *user, // todo fetch name from zultys api?? we need to cache this somewhere
			Extension: extension,
			Duration:  duration,
			DialPlan:  trunk,
			Caller:    caller,
			Callee:    callee,
		}

		calls = append(calls, entry)
		extensionTimestamps[extension] = callTime.Unix()
		// todo process the actual lines?

		/*err = SaveTimestampToFile(config.TimestampFile, oldestTimestamp)
		if err != nil {
			log.Error(err)
			return nil, err
		}*/

		err = SaveTimestampsToFile(config.TimestampFile, extensionTimestamps)
		if err != nil {
			log.Errorf("Error saving updated timestamps: %v", err)
			return nil, err
		}
	}

	return calls, nil
}

func processFtpData(content DataContent, config Config) ([]CallEntry, error) {
	var calls []CallEntry

	if strings.Contains(content.FilePath, "Calls By Extension") {
		lines := strings.Split(content.Content, "\n")

		file, err := LoadUsersFromFile(config.ZultysUsersFile)
		if err != nil {
			log.Errorf("could not find user? %s", err)
		}

		/*oldestTimestamp := time.Time{}
		// get timestamp
		timestmp, err := GetTimestampFromFile(config.TimestampFile)
		// if fails to load, skip and mark oldest as empty, otherwise
		if err != nil {
			log.Error("unable to get timestamp from file, saving blank to file?")
		}

		if !timestmp.IsZero() {
			oldestTimestamp = timestmp
		}*/

		// Load previously processed timestamps
		extensionTimestamps, err := LoadTimestampsFromFile(config.TimestampFile)
		if err != nil {
			log.Errorf("unable to get timestamps from file, assuming new: %v", err)
			extensionTimestamps = make(ExtensionTimestamps)
		}

		for _, line := range lines {
			line = strings.ReplaceAll(line, "\r", "")
			if line == "" || line == "\r" {
				log.Info("skipping blank line")
				continue
			}
			line = strings.ReplaceAll(line, "\"", "")
			values := strings.Split(line, ",")

			// example format: "222","3/14/2024","11:30:35 AM","00:01:02","Outbound","12502025183","11 Digit"
			//                   0        1           2            3          4           5            6

			dialplan := values[6]
			if strings.Contains(dialplan, "park") {
				log.Info("skipping over park entry")
				continue
			}

			direction := INTERNAL
			if values[4] == "Outbound" {
				direction = OUTBOUND
			} else if values[4] == "Inbound" {
				direction = INBOUND
			}

			if direction == INTERNAL {
				continue
			}

			// todo parse time & duration

			const dateTimeLayout = "1/2/2006 3:04:05 PM" // Combined layout for parsing both together

			// Combine the date and time strings
			dateTimeStr := fmt.Sprintf("%s %s", values[1], values[2])

			location, err := time.LoadLocation(config.TimestampRegion)
			if err != nil {
				log.Fatalf("Error loading location: %v\n", err)
			}

			// Parse the combined date and time string into a time.Time object
			callTime, err := time.ParseInLocation(dateTimeLayout, dateTimeStr, location)
			if err != nil {
				log.Errorf("Error parsing date and time: %v\n", err)
			}

			extension := values[0]

			// Check if the call is newer than the last processed for this extension
			lastProcessed, exists := extensionTimestamps[extension]
			if exists && callTime.Unix() <= lastProcessed {
				// This call has already been processed, skip it
				log.Warn("skipping call, due to being already processed?")
				continue
			}

			// keep track of oldest processed timestamp
			/*if oldestTimestamp.IsZero() || oldestTimestamp.Unix() < callTime.Unix() {
				oldestTimestamp = callTime
			} else if oldestTimestamp.Unix() > callTime.Unix() {
				log.Info("time has already passed")
				continue
			}*/
			// todo fix this time conversion
			// callTime = callTime.Add(-time.Hour * 2)

			durationSep := strings.Split(values[3], ":")
			duration, err := time.ParseDuration(durationSep[0] + "h" + durationSep[1] + "m" + durationSep[2] + "s")
			if err != nil {
				return nil, err
			}

			// fmt.Println(line)

			caller := values[0]
			callee := values[5]

			if direction == INBOUND {
				caller = values[5]
				callee = values[0]
			}

			user := FindEntryByExtension(file, values[0])
			if user == nil {
				log.Error("no user found")
				user = &ZultysUser{}
			}

			entry := CallEntry{
				Time:      callTime,
				Direction: direction,
				User:      *user, // todo fetch name from zultys api?? we need to cache this somewhere
				Extension: values[0],
				Duration:  duration,
				DialPlan:  dialplan,
				Caller:    caller,
				Callee:    callee,
			}

			calls = append(calls, entry)
			extensionTimestamps[extension] = callTime.Unix()
			// todo process the actual lines?
		}

		/*err = SaveTimestampToFile(config.TimestampFile, oldestTimestamp)
		if err != nil {
			log.Error(err)
			return nil, err
		}*/

		err = SaveTimestampsToFile(config.TimestampFile, extensionTimestamps)
		if err != nil {
			log.Errorf("Error saving updated timestamps: %v", err)
			return nil, err
		}
	}

	return calls, nil
}

// LoadTimestampsFromFile loads the timestamps for each extension from a JSON file.
func LoadTimestampsFromFile(filename string) (ExtensionTimestamps, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var timestamps ExtensionTimestamps
	err = json.Unmarshal(data, &timestamps)
	if err != nil {
		return nil, err
	}
	return timestamps, nil
}

// SaveTimestampsToFile saves the timestamps for each extension to a JSON file.
func SaveTimestampsToFile(filename string, timestamps ExtensionTimestamps) error {
	jsonData, err := json.Marshal(timestamps)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filename, jsonData, 0644)
}

type ExtensionTimestamps map[string]int64 // extension -> timestamp

var extensionTimestamps ExtensionTimestamps
