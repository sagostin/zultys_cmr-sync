package hubspot

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"time"
)

type CompanyProperties struct {
	City               string    `json:"city"`
	CreateDate         time.Time `json:"createdate"`
	HsTimestamp        time.Time `json:"hs_timestamp"`
	Domain             string    `json:"domain"`
	HsLastModifiedDate time.Time `json:"lastmodifieddate"`
	Industry           string    `json:"industry"`
	FullName           string    `json:"full_name"`
	Phone              string    `json:"phone"`
	State              string    `json:"state"`
}

type CompanySearchResponse struct {
	Total   int `json:"total"`
	Results []struct {
		Id         string            `json:"id"`
		Properties CompanyProperties `json:"properties"`
		CreatedAt  time.Time         `json:"createdAt"`
		UpdatedAt  time.Time         `json:"updatedAt"`
		Archived   bool              `json:"archived"`
	} `json:"results"`
	Paging struct {
		Next struct {
			After string `json:"after"`
			Link  string `json:"link"`
		} `json:"next"`
	} `json:"paging"`
}

// Companies client
type Companies struct {
	Client
}

// Companies constructor (from Client)
func (c Client) Companies() Companies {
	return Companies{
		Client: c,
	}
}

func (c Companies) SearchByPhone(phoneNumber string) (CompanySearchResponse, error) {
	filterQuery := FilterQuery{
		FilterGroups: []FilterGroup{},
		//Properties:   []string{"hs_searchable_calculated_phone"},
	}

	fGroup := FilterGroup{}
	fGroup.Filters = []Filter{}
	fGroup.Filters = append(fGroup.Filters, Filter{
		Value:        phoneNumber,
		PropertyName: "phone",
		Operator:     "EQ",
	})

	filterQuery.FilterGroups = append(filterQuery.FilterGroups, fGroup)

	j, _ := json.Marshal(filterQuery)
	log.Info(string(j))

	resp := CompanySearchResponse{}

	err := c.Client.Request("POST", "/crm/v3/objects/companies/search", &filterQuery, &resp)
	return resp, err
}
