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

type Account struct {
	Owner struct {
		Name  string `json:"name"`
		Id    string `json:"id"`
		Email string `json:"email"`
	} `json:"Owner"`
	CurrencySymbol    string      `json:"$currency_symbol"`
	FieldStates       interface{} `json:"$field_states"`
	AccountType       string      `json:"Account_Type"`
	SICCode           interface{} `json:"SIC_Code"`
	SharingPermission string      `json:"$sharing_permission"`
	LastActivityTime  time.Time   `json:"Last_Activity_Time"`
	Industry          string      `json:"Industry"`
	AccountSite       interface{} `json:"Account_Site"`
	ProcessFlow       bool        `json:"$process_flow"`
	BillingCountry    string      `json:"Billing_Country"`
	LockedForMe       bool        `json:"$locked_for_me"`
	Id                string      `json:"id"`
	Approval          struct {
		Delegate bool `json:"delegate"`
		Approve  bool `json:"approve"`
		Reject   bool `json:"reject"`
		Resubmit bool `json:"resubmit"`
	} `json:"$approval"`
	BillingStreet        string      `json:"Billing_Street"`
	CreatedTime          time.Time   `json:"Created_Time"`
	WizardConnectionPath interface{} `json:"$wizard_connection_path"`
	Editable             bool        `json:"$editable"`
	BillingCode          string      `json:"Billing_Code"`
	ShippingCity         interface{} `json:"Shipping_City"`
	ShippingCountry      interface{} `json:"Shipping_Country"`
	ShippingCode         interface{} `json:"Shipping_Code"`
	BillingCity          string      `json:"Billing_City"`
	CreatedBy            struct {
		Name  string `json:"name"`
		Id    string `json:"id"`
		Email string `json:"email"`
	} `json:"Created_By"`
	ZiaOwnerAssignment string      `json:"$zia_owner_assignment"`
	AnnualRevenue      int         `json:"Annual_Revenue"`
	ShippingStreet     interface{} `json:"Shipping_Street"`
	Ownership          string      `json:"Ownership"`
	Description        string      `json:"Description"`
	Rating             interface{} `json:"Rating"`
	ShippingState      interface{} `json:"Shipping_State"`
	ReviewProcess      struct {
		Approve  bool `json:"approve"`
		Reject   bool `json:"reject"`
		Resubmit bool `json:"resubmit"`
	} `json:"$review_process"`
	Website     string `json:"Website"`
	Employees   int    `json:"Employees"`
	RecordImage string `json:"Record_Image"`
	ModifiedBy  struct {
		Name  string `json:"name"`
		Id    string `json:"id"`
		Email string `json:"email"`
	} `json:"Modified_By"`
	Review        interface{}   `json:"$review"`
	Phone         string        `json:"Phone"`
	AccountName   string        `json:"Account_Name"`
	ZiaVisions    interface{}   `json:"$zia_visions"`
	AccountNumber string        `json:"Account_Number"`
	TickerSymbol  interface{}   `json:"Ticker_Symbol"`
	ModifiedTime  time.Time     `json:"Modified_Time"`
	RecordStatusS string        `json:"Record_Status__s"`
	Orchestration bool          `json:"$orchestration"`
	ParentAccount interface{}   `json:"Parent_Account"`
	InMerge       bool          `json:"$in_merge"`
	LockedS       bool          `json:"Locked__s"`
	BillingState  string        `json:"Billing_State"`
	Tag           []interface{} `json:"Tag"`
	Fax           interface{}   `json:"Fax"`
	ApprovalState string        `json:"$approval_state"`
	Pathfinder    bool          `json:"$pathfinder"`
}

type AccountSearchResponse struct {
	Data []Account `json:"data"`
	Info struct {
		PerPage     int    `json:"per_page"`
		Count       int    `json:"count"`
		SortBy      string `json:"sort_by"`
		Page        int    `json:"page"`
		SortOrder   string `json:"sort_order"`
		MoreRecords bool   `json:"more_records"`
	} `json:"info"`
}

/*type AccountResponse struct {
	Fields []struct {
		AssociatedModule interface{} `json:"associated_module"`
		Webhook          bool        `json:"webhook"`
		OperationType    struct {
			WebUpdate bool `json:"web_update"`
			ApiCreate bool `json:"api_create"`
			WebCreate bool `json:"web_create"`
			ApiUpdate bool `json:"api_update"`
		} `json:"operation_type"`
		ColourCodeEnabledBySystem bool        `json:"colour_code_enabled_by_system"`
		FieldLabel                string      `json:"field_label"`
		Tooltip                   interface{} `json:"tooltip"`
		Type                      string      `json:"type"`
		FieldReadOnly             bool        `json:"field_read_only"`
		CustomizableProperties    []string    `json:"customizable_properties"`
		DisplayLabel              string      `json:"display_label"`
		ReadOnly                  bool        `json:"read_only"`
		AssociationDetails        interface{} `json:"association_details"`
		BusinesscardSupported     bool        `json:"businesscard_supported"`
		MultiModuleLookup         struct {
		} `json:"multi_module_lookup"`
		Id          string      `json:"id"`
		CreatedTime interface{} `json:"created_time"`
		Filterable  bool        `json:"filterable"`
		Visible     bool        `json:"visible"`
		Profiles    []struct {
			PermissionType string `json:"permission_type"`
			Name           string `json:"name"`
			Id             string `json:"id"`
		} `json:"profiles"`
		ViewType struct {
			View        bool `json:"view"`
			Edit        bool `json:"edit"`
			QuickCreate bool `json:"quick_create"`
			Create      bool `json:"create"`
		} `json:"view_type"`
		Separator  bool        `json:"separator"`
		Searchable bool        `json:"searchable"`
		External   interface{} `json:"external"`
		ApiName    string      `json:"api_name"`
		Unique     struct {
		} `json:"unique"`
		EnableColourCode bool `json:"enable_colour_code"`
		PickListValues   []struct {
			DisplayValue   string  `json:"display_value"`
			SequenceNumber int     `json:"sequence_number"`
			ReferenceValue string  `json:"reference_value"`
			ColourCode     *string `json:"colour_code"`
			ActualValue    string  `json:"actual_value"`
			Id             string  `json:"id"`
			Type           string  `json:"type"`
		} `json:"pick_list_values"`
		SystemMandatory bool        `json:"system_mandatory"`
		VirtualField    bool        `json:"virtual_field"`
		JsonType        string      `json:"json_type"`
		Crypt           interface{} `json:"crypt"`
		CreatedSource   string      `json:"created_source"`
		DisplayType     int         `json:"display_type"`
		UiType          int         `json:"ui_type"`
		ModifiedTime    *time.Time  `json:"modified_time"`
		EmailParser     struct {
			FieldsUpdateSupported     bool `json:"fields_update_supported"`
			RecordOperationsSupported bool `json:"record_operations_supported"`
		} `json:"email_parser"`
		Currency struct {
			RoundingOption string `json:"rounding_option,omitempty"`
			Precision      int    `json:"precision,omitempty"`
		} `json:"currency"`
		CustomField bool `json:"custom_field"`
		Lookup      struct {
			DisplayLabel               string `json:"display_label,omitempty"`
			RevalidateFilterDuringEdit bool   `json:"revalidate_filter_during_edit,omitempty"`
			ApiName                    string `json:"api_name,omitempty"`
			Module                     struct {
				ApiName string `json:"api_name"`
				Id      string `json:"id"`
			} `json:"module,omitempty"`
			Id           string `json:"id,omitempty"`
			QueryDetails struct {
			} `json:"query_details,omitempty"`
		} `json:"lookup"`
		RollupSummary struct {
		} `json:"rollup_summary"`
		Length                        int         `json:"length"`
		DisplayField                  bool        `json:"display_field"`
		PickListValuesSortedLexically bool        `json:"pick_list_values_sorted_lexically"`
		Sortable                      bool        `json:"sortable"`
		GlobalPicklist                interface{} `json:"global_picklist"`
		HistoryTracking               interface{} `json:"history_tracking"`
		DataType                      string      `json:"data_type"`
		Formula                       struct {
		} `json:"formula"`
		DecimalPlace      *int `json:"decimal_place"`
		MassUpdate        bool `json:"mass_update"`
		Multiselectlookup struct {
		} `json:"multiselectlookup"`
		AutoNumber struct {
		} `json:"auto_number"`
		BlueprintSupported  bool   `json:"blueprint_supported,omitempty"`
		QuickSequenceNumber string `json:"quick_sequence_number,omitempty"`
		Textarea            struct {
			Type string `json:"type"`
		} `json:"textarea,omitempty"`
	} `json:"fields"`
}*/

// endpoint: https://crm.zohocloud.ca/crm/v6/Contacts?fields=
// Owner,Rating,Account_Name,Phone,Account_Site,Fax,Parent_Account,
// Website,Account_Number,Ticker_Symbol,Account_Type,Ownership,
// Industry,Employees,Annual_Revenue,SIC_Code,Tag,Created_By,
// Modified_By,Created_Time,Modified_Time,Last_Activity_Time,id,
// Change_Log_Time__s,Enrich_Status__s,Last_Enriched_Time__s,Record_Status__s,
// Locked__s,Billing_Street,Shipping_Street,Billing_City,Shipping_City,
// Billing_State,Shipping_State,Billing_Code,Shipping_Code,Billing_Country,
// Shipping_Country,Description,Record_Image

func (c *Client) FindAccountByPhone(phone string) (AccountSearchResponse, error) {
	var accountResponse AccountSearchResponse

	// Construct the GET request to fetch contacts
	req, err := http.NewRequest("GET", "https://"+c.Endpoints.CrmApi+"/crm/v6/Accounts/search?phone="+phone+"", nil)
	if err != nil {
		return accountResponse, fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Add("Authorization", "Zoho-oauthtoken "+c.Auth.AccessToken)

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return accountResponse, fmt.Errorf("failed to send request: %v", err)
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
		return accountResponse, fmt.Errorf("failed to read response: %v", err)
	}

	// Unmarshal the response
	if err := json.Unmarshal(body, &accountResponse); err != nil {
		return accountResponse, fmt.Errorf("failed to unmarshal response: %v", err)
	}

	return accountResponse, nil
}
