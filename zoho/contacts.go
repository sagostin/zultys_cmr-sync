package zoho

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

type Contact struct {
	Owner struct {
		Name  string `json:"name"`
		Id    string `json:"id"`
		Email string `json:"email"`
	} `json:"Owner"`
	Email             string      `json:"Email"`
	CurrencySymbol    string      `json:"$currency_symbol"`
	FieldStates       interface{} `json:"$field_states"`
	OtherPhone        interface{} `json:"Other_Phone"`
	MailingState      string      `json:"Mailing_State"`
	OtherState        interface{} `json:"Other_State"`
	SharingPermission string      `json:"$sharing_permission"`
	OtherCountry      interface{} `json:"Other_Country"`
	LastActivityTime  time.Time   `json:"Last_Activity_Time"`
	Department        string      `json:"Department"`
	UnsubscribedMode  interface{} `json:"Unsubscribed_Mode"`
	ProcessFlow       bool        `json:"$process_flow"`
	Assistant         interface{} `json:"Assistant"`
	MailingCountry    string      `json:"Mailing_Country"`
	LockedForMe       bool        `json:"$locked_for_me"`
	Id                string      `json:"id"`
	ReportingTo       interface{} `json:"Reporting_To"`
	Approval          struct {
		Delegate bool `json:"delegate"`
		Approve  bool `json:"approve"`
		Reject   bool `json:"reject"`
		Resubmit bool `json:"resubmit"`
	} `json:"$approval"`
	OtherCity            interface{} `json:"Other_City"`
	CreatedTime          time.Time   `json:"Created_Time"`
	WizardConnectionPath interface{} `json:"$wizard_connection_path"`
	Editable             bool        `json:"$editable"`
	HomePhone            interface{} `json:"Home_Phone"`
	CreatedBy            struct {
		Name  string `json:"name"`
		Id    string `json:"id"`
		Email string `json:"email"`
	} `json:"Created_By"`
	ZiaOwnerAssignment string      `json:"$zia_owner_assignment"`
	SecondaryEmail     interface{} `json:"Secondary_Email"`
	Description        interface{} `json:"Description"`
	VendorName         interface{} `json:"Vendor_Name"`
	MailingZip         string      `json:"Mailing_Zip"`
	ReviewProcess      struct {
		Approve  bool `json:"approve"`
		Reject   bool `json:"reject"`
		Resubmit bool `json:"resubmit"`
	} `json:"$review_process"`
	Twitter       string      `json:"Twitter"`
	OtherZip      interface{} `json:"Other_Zip"`
	MailingStreet string      `json:"Mailing_Street"`
	Salutation    interface{} `json:"Salutation"`
	FirstName     string      `json:"First_Name"`
	FullName      string      `json:"Full_Name"`
	AsstPhone     interface{} `json:"Asst_Phone"`
	RecordImage   string      `json:"Record_Image"`
	ModifiedBy    struct {
		Name  string `json:"name"`
		Id    string `json:"id"`
		Email string `json:"email"`
	} `json:"Modified_By"`
	Review      interface{} `json:"$review"`
	SkypeID     string      `json:"Skype_ID"`
	Phone       string      `json:"Phone"`
	AccountName struct {
		Name string `json:"name"`
		Id   string `json:"id"`
	} `json:"Account_Name"`
	EmailOptOut      bool          `json:"Email_Opt_Out"`
	ZiaVisions       interface{}   `json:"$zia_visions"`
	ModifiedTime     time.Time     `json:"Modified_Time"`
	DateOfBirth      interface{}   `json:"Date_of_Birth"`
	MailingCity      string        `json:"Mailing_City"`
	UnsubscribedTime interface{}   `json:"Unsubscribed_Time"`
	Title            string        `json:"Title"`
	OtherStreet      interface{}   `json:"Other_Street"`
	Mobile           string        `json:"Mobile"`
	RecordStatusS    string        `json:"Record_Status__s"`
	Orchestration    bool          `json:"$orchestration"`
	LastName         string        `json:"Last_Name"`
	InMerge          bool          `json:"$in_merge"`
	LockedS          bool          `json:"Locked__s"`
	LeadSource       string        `json:"Lead_Source"`
	Tag              []interface{} `json:"Tag"`
	Fax              interface{}   `json:"Fax"`
	ApprovalState    string        `json:"$approval_state"`
	Pathfinder       bool          `json:"$pathfinder"`
}

type ContactSearchResponse struct {
	Data []Contact `json:"data"`
	Info struct {
		PerPage     int    `json:"per_page"`
		Count       int    `json:"count"`
		SortBy      string `json:"sort_by"`
		Page        int    `json:"page"`
		SortOrder   string `json:"sort_order"`
		MoreRecords bool   `json:"more_records"`
	} `json:"info"`
}

// endpoint: https://crm.zohocloud.ca/crm/v6/Contacts?fields=
// Owner,First_Name,Salutation,Last_Name,Full_Name,Account_Name,Vendor_Name,
// Email,Title,Department,Phone,Home_Phone,Other_Phone,Fax,Mobile,Date_of_Birth,
// Tag,Assistant,Reporting_To,Email_Opt_Out,Created_By,Skype_ID,Modified_By,
// Created_Time,Modified_Time,Secondary_Email,Last_Activity_Time,Twitter,id,
// Change_Log_Time__s,Record_Status__s,Unsubscribed_Mode,Unsubscribed_Time,
// Enrich_Status__s,Last_Enriched_Time__s,Locked__s,Mailing_Street,Other_Street,
// Mailing_City,Other_City,Mailing_State,Other_State,Mailing_Zip,Other_Zip,
// Mailing_Country,Other_Country,Description,Record_Image

func (c *Client) FindContactByPhone(phone string) (ContactSearchResponse, error) {
	var contactResponse ContactSearchResponse

	// Construct the GET request to fetch contacts
	req, err := http.NewRequest("GET", "https://"+c.Endpoints.CrmApi+"/crm/v6/Contacts/search?phone="+phone+"", nil)
	if err != nil {
		return contactResponse, fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Add("Authorization", "Zoho-oauthtoken "+c.Auth.AccessToken)

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return contactResponse, fmt.Errorf("failed to send request: %v", err)
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
		return contactResponse, fmt.Errorf("failed to read response: %v", err)
	}

	// Unmarshal the response
	if err := json.Unmarshal(body, &contactResponse); err != nil {
		return contactResponse, fmt.Errorf("failed to unmarshal response: %v", err)
	}

	return contactResponse, nil
}

/*func (c *Client) FetchContacts() (ContactResponse, error) {
	var contactResponse ContactResponse

	// Construct the GET request to fetch contacts
	req, err := http.NewRequest("GET", "https://"+c.Endpoints.CrmApi+"/crm/v6/Contacts", nil)
	if err != nil {
		return contactResponse, fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Add("Authorization", "Zoho-oauthtoken "+c.Auth.AccessToken)

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return contactResponse, fmt.Errorf("failed to send request: %v", err)
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
		return contactResponse, fmt.Errorf("failed to read response: %v", err)
	}

	// Unmarshal the response
	if err := json.Unmarshal(body, &contactResponse); err != nil {
		return contactResponse, fmt.Errorf("failed to unmarshal response: %v", err)
	}

	return contactResponse, nil
}*/
