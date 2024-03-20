package zoho

import "time"

type AccountResponse struct {
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
}

// endpoint: https://crm.zohocloud.ca/crm/v6/Contacts?fields=
// Owner,Rating,Account_Name,Phone,Account_Site,Fax,Parent_Account,
// Website,Account_Number,Ticker_Symbol,Account_Type,Ownership,
// Industry,Employees,Annual_Revenue,SIC_Code,Tag,Created_By,
// Modified_By,Created_Time,Modified_Time,Last_Activity_Time,id,
// Change_Log_Time__s,Enrich_Status__s,Last_Enriched_Time__s,Record_Status__s,
// Locked__s,Billing_Street,Shipping_Street,Billing_City,Shipping_City,
// Billing_State,Shipping_State,Billing_Code,Shipping_Code,Billing_Country,
// Shipping_Country,Description,Record_Image
