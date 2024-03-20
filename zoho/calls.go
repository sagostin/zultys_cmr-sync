package zoho

import "time"

type Call struct {
	CallDuration string `json:"Call_Duration"`
	Owner        struct {
		Name  string `json:"name"`
		Id    string `json:"id"`
		Email string `json:"email"`
	} `json:"Owner"`
	Description string `json:"Description"`
	CallerID    string `json:"Caller_ID"`
	CTIEntry    bool   `json:"CTI_Entry"`
	CallAgenda  string `json:"Call_Agenda"`
	ModifiedBy  struct {
		Name  string `json:"name"`
		Id    string `json:"id"`
		Email string `json:"email"`
	} `json:"Modified_By"`
	CallPurpose        string    `json:"Call_Purpose"`
	Id                 string    `json:"id"`
	WhoId              string    `json:"Who_Id"`
	OutgoingCallStatus string    `json:"Outgoing_Call_Status"`
	ModifiedTime       time.Time `json:"Modified_Time"`
	Reminder           time.Time `json:"Reminder"`
	//VoiceRecordingS    interface{} `json:"Voice_Recording__s"`
	CreatedTime    time.Time `json:"Created_Time"`
	CallStartTime  time.Time `json:"Call_Start_Time"`
	Subject        string    `json:"Subject"`
	SeModule       string    `json:"$se_module"`
	CallType       string    `json:"Call_Type"`
	ScheduledInCRM string    `json:"Scheduled_In_CRM"`
	CallResult     string    `json:"Call_Result"`
	WhatId         struct {
		Name string `json:"name"`
		Id   string `json:"id"`
	} `json:"What_Id"`
	CallDurationInSeconds interface{} `json:"Call_Duration_in_seconds"`
	CreatedBy             struct {
		Name  string `json:"name"`
		Id    string `json:"id"`
		Email string `json:"email"`
	} `json:"Created_By"`
	Tag           string `json:"Tag"`
	DialledNumber string `json:"Dialled_Number"`
}

type CallsResponse struct {
	Data []Call `json:"data"`
	Info struct {
		PerPage           int         `json:"per_page"`
		NextPageToken     interface{} `json:"next_page_token"`
		Count             int         `json:"count"`
		SortBy            string      `json:"sort_by"`
		Page              int         `json:"page"`
		PreviousPageToken interface{} `json:"previous_page_token"`
		PageTokenExpiry   interface{} `json:"page_token_expiry"`
		SortOrder         string      `json:"sort_order"`
		MoreRecords       bool        `json:"more_records"`
	} `json:"info"`
}

// https://crm.zohocloud.ca/crm/v6/Calls?fields=
// Who_Id,What_Id,Call_Type,Outgoing_Call_Status,Call_Start_Time,Call_Duration,Owner,Dialled_Number,Subject,Caller_ID,
// Created_By,Modified_By,Created_Time,Modified_Time,Reminder,Voice_Recording__s,Call_Duration_in_seconds,Scheduled_In_CRM,
// CTI_Entry,Tag,id,Call_Purpose,Call_Agenda,Call_Result,Description
