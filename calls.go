package hubspot

import (
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
	CreateDate         time.Time `json:"createdate,omitempty"`
	HsCallBody         string    `json:"hs_call_body,omitempty"`
	HsCallDuration     string    `json:"hs_call_duration,omitempty"`
	HsCallFromNumber   string    `json:"hs_call_from_number,omitempty"`
	HsCallRecordingUrl string    `json:"hs_call_recording_url,omitempty"`
	HsCallStatus       string    `json:"hs_call_status,omitempty"`
	HsCallTitle        string    `json:"hs_call_title,omitempty"`
	HsCallToNumber     string    `json:"hs_call_to_number,omitempty"`
	HsLastModifiedDate time.Time `json:"hs_lastmodifieddate,omitempty"`
	HsTimestamp        time.Time `json:"hs_timestamp,omitempty"`
	HubspotOwnerId     string    `json:"hubspot_owner_id,omitempty"`
}

type MultiCallProperties struct {
	HsCreateDate       time.Time `json:"hs_createdate"`
	HsLastModifiedDate time.Time `json:"hs_lastmodifieddate"`
	HsObjectId         string    `json:"hs_object_id"`
}

type CallsResponse struct {
	Total   int `json:"total"`
	Results []struct {
		Id         string         `json:"id"`
		Properties CallProperties `json:"properties"`
		CreatedAt  time.Time      `json:"createdAt"`
		UpdatedAt  time.Time      `json:"updatedAt"`
		Archived   bool           `json:"archived"`
	} `json:"results"`
	Paging struct {
		Next struct {
			After string `json:"after"`
		} `json:"next"`
	} `json:"paging"`
}

func (c Calls) GetCalls() (CallsResponse, error) {
	//GET /crm/v3/objects/calls
	resp := CallsResponse{}

	err := c.Client.Request("GET", "/crm/v3/objects/calls", nil, &resp)
	return resp, err
}

func (c Calls) GetAllCalls(t time.Time) ([]CallProperties, error) {

	var n, l int
	n = 0
	l = 0

	var inc int
	inc = 25

	filterQuery := FilterQuery{
		Properties: []string{"hs_call_from_number", "hs_createdate", "hs_call_callee_object_type_id", "hs_call_callee_object_id"},
		Sorts: []FilterSort{{
			PropertyName: "hs_lastmodifieddate",
			Direction:    "DESCENDING",
		}},
		Limit: inc,
		After: n,
	}

	var props []CallProperties

	resp := CallsResponse{}
	err := c.Client.Request("POST", "/crm/v3/objects/calls/search", &filterQuery, &resp)
	if err != nil {
		log.Error(err)
	}
	l = resp.Total
	n = inc
	for _, r := range resp.Results {
		props = append(props, r.Properties)
	}

	log.Infof("Total Calls: %v", resp.Total)

	for n < l {
		resp = CallsResponse{}
		err = c.Client.Request("POST", "/crm/v3/objects/calls/search", &filterQuery, &resp)
		if err != nil {
			log.Error(err)
		}
		n += inc
		for _, r := range resp.Results {
			props = append(props, r.Properties)
		}
	}

	/*fGroup := FilterGroup{}
	fGroup.Filters = []Filter{}
	fGroup.Filters = append(fGroup.Filters, Filter{
		Value:        t.Format(time.RFC3339),
		PropertyName: "createdate",
		Operator:     "GTE",
	})

	filterQuery.FilterGroups = append(filterQuery.FilterGroups, fGroup)
	marshal, err := json.Marshal(filterQuery)
	if err != nil {
		log.Fatalf("%s", "Error parsing json.")
	}
	log.Infof("%s", marshal)*/

	//GET /crm/v3/objects/calls
	return props, err
}

func (c Calls) GetCall(objId string) (CallProperties, error) {
	/*filterQuery := FilterQuery{
		Properties: []string{"hs_call_from_number", "hs_call_title"},
	}

	fGroup := FilterGroup{}
	fGroup.Filters = []Filter{}
	fGroup.Filters = append(fGroup.Filters, Filter{
		Value:        objId,
		PropertyName: "hs_object_id",
		Operator:     "EQ",
	})

	filterQuery.FilterGroups = append(filterQuery.FilterGroups, fGroup)

	marshal, err := json.Marshal(filterQuery)
	if err != nil {
		log.Fatalf("%s", "Error parsing json.")
	}
	log.Infof("%s", marshal)*/

	resp := CallProperties{}

	err := c.Client.Request("GET", "/crm/v3/objects/calls/"+objId, nil, &resp)
	return resp, err
}
