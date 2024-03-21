package zoho

import (
	"bytes"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

type Owner struct {
	Name  string `json:"name,omitempty"`
	Id    string `json:"id"`
	Email string `json:"email,omitempty"`
}

type ItemLink struct {
	Name string `json:"name,omitempty"`
	Id   string `json:"id"`
}

type Call struct {
	CallDuration       string    `json:"Call_Duration,omitempty"`
	Owner              *Owner    `json:"Owner"`
	Description        string    `json:"Description,omitempty"`
	CallerID           string    `json:"Caller_ID,omitempty"`
	CTIEntry           bool      `json:"CTI_Entry,omitempty"`
	CallAgenda         string    `json:"Call_Agenda,omitempty"`
	ModifiedBy         *Owner    `json:"Modified_By"`
	CallPurpose        string    `json:"Call_Purpose,omitempty"`
	Id                 string    `json:"id,omitempty"`
	WhoId              *ItemLink `json:"Who_Id"`
	OutgoingCallStatus string    `json:"Outgoing_Call_Status,omitempty"`
	ModifiedTime       time.Time `json:"Modified_Time,omitempty"`
	Reminder           time.Time `json:"Reminder,omitempty"`
	//VoiceRecordingS    interface{} `json:"Voice_Recording__s"`
	CreatedTime           time.Time   `json:"Created_Time,omitempty"`
	CallStartTime         time.Time   `json:"Call_Start_Time,omitempty"`
	Subject               string      `json:"Subject,omitempty"`
	SeModule              string      `json:"$se_module,omitempty"`
	CallType              string      `json:"Call_Type,omitempty"`
	ScheduledInCRM        string      `json:"Scheduled_In_CRM,omitempty"`
	CallResult            string      `json:"Call_Result,omitempty"`
	WhatId                *ItemLink   `json:"What_Id"`
	CallDurationInSeconds interface{} `json:"Call_Duration_in_seconds,omitempty"`
	CreatedBy             *Owner      `json:"Created_By"`
	Tag                   string      `json:"Tag,omitempty"`
	DialledNumber         string      `json:"Dialled_Number,omitempty"`
}

type CallsResponse struct {
	Data []Call `json:"data"`
	Info struct {
		PerPage           int         `json:"per_page,omitempty"`
		NextPageToken     interface{} `json:"next_page_token,omitempty"`
		Count             int         `json:"count,omitempty"`
		SortBy            string      `json:"sort_by,omitempty"`
		Page              int         `json:"page,omitempty"`
		PreviousPageToken interface{} `json:"previous_page_token,omitempty"`
		PageTokenExpiry   interface{} `json:"page_token_expiry,omitempty"`
		SortOrder         string      `json:"sort_order,omitempty"`
		MoreRecords       bool        `json:"more_records,omitempty"`
	} `json:"info,omitempty"`
}

// https://crm.zohocloud.ca/crm/v6/Calls?fields=
// Who_Id,What_Id,Call_Type,Outgoing_Call_Status,Call_Start_Time,Call_Duration,Owner,Dialled_Number,Subject,Caller_ID,
// Created_By,Modified_By,Created_Time,Modified_Time,Reminder,Voice_Recording__s,Call_Duration_in_seconds,Scheduled_In_CRM,
// CTI_Entry,Tag,id,Call_Purpose,Call_Agenda,Call_Result,Description

func (c *Client) CreateCall(call []Call) error {
	callData, err := json.Marshal(struct {
		Data []Call `json:"data"`
	}{
		Data: call,
	})
	if err != nil {
		return fmt.Errorf("failed to marshal call data: %v", err)
	}

	req, err := http.NewRequest("POST", "https://"+c.Endpoints.CrmApi+"/crm/v6/Calls", bytes.NewBuffer(callData))
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Add("Authorization", "Zoho-oauthtoken "+c.Auth.AccessToken)

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %v", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Error(err)
		}
	}(resp.Body)

	// Read the response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %v", err)
	}
	fmt.Println(string(body))

	return nil
}
