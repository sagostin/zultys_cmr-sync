package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
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
	Name      string
	Extension string
	Duration  time.Duration
	DialPlan  string
	Caller    string
	Callee    string
}

func processData(content DataContent) ([]CallEntry, error) {
	if content.Type == DataTypeFTP {
		data, err := processFtpData(content)
		if err != nil {
			return data, err
		}

		for _, d := range data {
			log.Infof("%s", d)
		}
	} else if content.Type == DataTypeSMDR {
		// todo handle smdr data format?
	}
	return nil, nil
}

func processFtpData(content DataContent) ([]CallEntry, error) {
	var calls []CallEntry

	if strings.Contains(content.FilePath, "Calls By Extension") {
		lines := strings.Split(content.Content, "\n")

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

			// Parse the combined date and time string into a time.Time object
			dateTime, err := time.Parse(dateTimeLayout, dateTimeStr)
			if err != nil {
				log.Errorf("Error parsing date and time: %v\n", err)
			}

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

			entry := CallEntry{
				Time:      dateTime,
				Direction: direction,
				Name:      "PLACEHOLDER", // todo fetch name from zultys api?? we need to cache this somewhere
				Extension: values[0],
				Duration:  duration,
				DialPlan:  dialplan,
				Caller:    caller,
				Callee:    callee,
			}

			calls = append(calls, entry)
			// todo process the actual lines?
		}
	}

	return calls, nil
}
