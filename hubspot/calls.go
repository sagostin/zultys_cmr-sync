package hubspot

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"time"
)

// Calls client
type Calls struct {
	Client
}

// Calls constructor (from Client)
func (c Client) Calls() Calls {
	return Calls{
		Client: c,
	}
}

type CallProperties struct {
	HsCallBody               string    `json:"hs_call_body,omitempty"`
	HsCallDuration           string    `json:"hs_call_duration,omitempty"`
	HsCallFromNumber         string    `json:"hs_call_from_number,omitempty"`
	HsCallRecordingUrl       string    `json:"hs_call_recording_url,omitempty"`
	HsCallStatus             string    `json:"hs_call_status,omitempty"`
	HsCallTitle              string    `json:"hs_call_title,omitempty"`
	HsCallDirection          string    `json:"hs_call_direction"`
	HsCallToNumber           string    `json:"hs_call_to_number,omitempty"`
	HsLastModifiedDate       time.Time `json:"hs_lastmodifieddate,omitempty"`
	HsTimestamp              time.Time `json:"hs_timestamp,omitempty"`
	HubspotOwnerId           string    `json:"hubspot_owner_id,omitempty"`
	HsAttachmentIds          string    `json:"hs_attachment_ids,omitempty"`
	HsCallCalleeObjectId     string    `json:"hs_call_callee_object_id,omitempty"`
	HsCallCalleeObjectTypeId string    `json:"hs_call_callee_object_type_id,omitempty"`
	HsCallDisposition        string    `json:"hs_call_disposition"`
}

type MultiCallProperties struct {
	HsCreateDate       time.Time `json:"hs_createdate"`
	HsLastModifiedDate time.Time `json:"hs_lastmodifieddate"`
	HsObjectId         string    `json:"hs_object_id"`
}

type CallResult struct {
	Id         string         `json:"id"`
	Properties CallProperties `json:"properties"`
	CreatedAt  time.Time      `json:"createdAt"`
	UpdatedAt  time.Time      `json:"updatedAt"`
	Archived   bool           `json:"archived"`
}

type CallsResponse struct {
	Total   int          `json:"total"`
	Results []CallResult `json:"results"`
	Paging  struct {
		Next struct {
			After string `json:"after"`
		} `json:"next"`
	} `json:"paging"`
}

type CallAssociateResponse struct {
	Properties CallProperties `json:"properties"`
}

func (c Calls) GetCalls() (CallsResponse, error) {
	//GET /crm/v3/objects/calls
	resp := CallsResponse{}

	err := c.Client.Request("GET", "/crm/v3/objects/calls", nil, &resp)
	return resp, err
}

func (c Calls) CreateCall(properties CallProperties) (*CallResult, error) {
	var call CallResult

	err := c.Client.Request("POST", "/crm/v3/objects/calls", CallAssociateResponse{Properties: properties}, &call)
	if err != nil {
		log.Error("error creating call record in hubspot")
		return nil, err
	}

	return &call, nil
}

/*func (c Calls) GetRecentCalls(duration time.Duration) ([]CallResult, error) {
	var afterCallX = 0
	var itemLimit = 25

	var notFinished = true
	var props []CallResult

	for notFinished {
		filterQuery := FilterQuery{
			Properties: []string{"hs_call_body", "hs_call_duration",
				"hs_call_from_number", "hs_call_recording_url",
				"hs_call_status", "hs_call_title", "hs_lastmodifieddate",
				"hs_timestamp", "hubspot_owner_id", "hs_call_to_number",
				"hs_attachment_ids", "hs_call_callee_object_id",
				"hs_call_direction",
				"hs_call_callee_object_type_id"},
			Sorts: []FilterSort{{
				PropertyName: "hs_timestamp",
				Direction:    "DESCENDING",
			}},
			Limit: itemLimit,
			After: afterCallX,
		}

		var resp = CallsResponse{}
		err := c.Client.Request("POST", "/crm/v3/objects/calls/search", &filterQuery, &resp)
		if err != nil {
			log.Error(err)
			return nil, err
		}
		for n, r := range resp.Results {

			log.Errorf("%s", r.Properties.HsTimestamp)
			if n == (itemLimit-1) && r.Properties.HsTimestamp.Before(time.Now().Add(-duration)) {
				//log.Infof("%s", time.Now().Add(-duration))
				notFinished = false
				break
			} else if n == (itemLimit-1) && r.Properties.HsTimestamp.After(time.Now().Add(-duration)) /* && resp.Total > afterCallX  {
				props = append(props, r)
				afterCallX = afterCallX + itemLimit
				log.Warnf("%v", afterCallX)
				continue
			} else if r.Properties.HsTimestamp.After(time.Now().Add(-duration)) {
				props = append(props, r)
				continue
			}
		}

		if !notFinished {
			break
		}

		time.Sleep(2 * time.Second)
	}

	return props, nil

}*/

func (c Calls) GetAllCalls() ([]CallResult, error) {
	var afterCallX = 0
	var itemLimit = 25

	var notFinished = true
	var props []CallResult

	for notFinished {
		filterQuery := FilterQuery{
			Properties: []string{"hs_call_body", "hs_call_duration",
				"hs_call_from_number", "hs_call_recording_url",
				"hs_call_status", "hs_call_title", "hs_lastmodifieddate",
				"hs_timestamp", "hubspot_owner_id", "hs_call_to_number",
				"hs_attachment_ids", "hs_call_callee_object_id",
				"hs_call_direction",
				"hs_call_callee_object_type_id", "hs_calculated_phone_number"},
			Sorts: []FilterSort{{
				PropertyName: "hs_timestamp",
				Direction:    "DESCENDING",
			}},
			Limit: itemLimit,
			After: afterCallX,
		}

		marshal, err := json.Marshal(filterQuery)
		if err != nil {
			return nil, err
		}

		log.Warn(string(marshal))

		var resp = CallsResponse{}
		err = c.Client.Request("POST", "/crm/v3/objects/calls/search", &filterQuery, &resp)
		if err != nil {
			log.Error(err)
			return nil, err
		}
		for n, r := range resp.Results {
			log.Errorf("%s", r.Properties.HsTimestamp)
			if n == (itemLimit-1) && (resp.Total > afterCallX+itemLimit) {
				afterCallX = afterCallX + itemLimit
				log.Warnf("%v", afterCallX)
				continue
			}
			log.Warnf("Position: %s , Total Pos: %s , Total Calls: %s", n, afterCallX+n, resp.Total)
			props = append(props, r)
			if resp.Total == afterCallX+n+1 {
				notFinished = false
				break
			}
			continue
		}

		if !notFinished {
			break
		}

		time.Sleep(2 * time.Second)
	}

	return props, nil

}

func (c Calls) GetCall(objId string) (CallResult, error) {
	resp := CallResult{}

	err := c.Client.Request("GET", "/crm/v3/objects/calls/"+objId, nil, &resp)
	return resp, err
}

func (c Calls) AssociateCallCompany(callResult CallResult, companyID string, label int) error {
	resp := CallAssociateResponse{}

	var assProps []AssociationProps
	assProps = append(assProps, AssociationProps{
		AssociationTypeId:   label,
		AssociationCategory: "HUBSPOT_DEFINED",
	})

	//log.Infof("%s", "/crm/v4/objects/calls/"+callResult.Id+"/associations/contact/"+contactResult.Id)

	err := c.Request("PUT", "/crm/v4/objects/calls/"+callResult.Id+"/associations/company/"+companyID, &assProps, &resp)
	if err != nil {
		return err
	}
	return nil
}

func (c Calls) AssociateCallContact(callResult CallResult, contactResult ContactResult, label int) error {
	resp := CallAssociateResponse{}

	var assProps []AssociationProps
	assProps = append(assProps, AssociationProps{
		AssociationTypeId:   label,
		AssociationCategory: "HUBSPOT_DEFINED",
	})

	//log.Infof("%s", "/crm/v4/objects/calls/"+callResult.Id+"/associations/contact/"+contactResult.Id)

	err := c.Request("PUT", "/crm/v4/objects/calls/"+callResult.Id+"/associations/contact/"+contactResult.Id, &assProps, &resp)
	if err != nil {
		return err
	}
	return nil
}
