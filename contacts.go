package hubspot

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"time"
)

// Contacts client
type Contacts struct {
	Client
}

// Contacts constructor (from Client)
func (c Client) Contacts() Contacts {
	return Contacts{
		Client: c,
	}
}

type ContactProperties struct {
	Company          string    `json:"company"`
	CreateDate       time.Time `json:"createdate"`
	Email            string    `json:"email"`
	Firstname        string    `json:"firstname"`
	LastModifiedDate time.Time `json:"lastmodifieddate"`
	Lastname         string    `json:"lastname"`
	Phone            string    `json:"phone"`
	Website          string    `json:"website"`
}

type ContactsSearchResponse struct {
	Total   int `json:"total"`
	Results []struct {
		Id         string            `json:"id"`
		Properties ContactProperties `json:"properties"`
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

func (c Contacts) SearchByPhone(phoneNumber string) (ContactsSearchResponse, error) {
	filterQuery := FilterQuery{
		FilterGroups: []FilterGroup{},
	}

	fGroup := FilterGroup{}
	fGroup.Filters = []Filter{}
	fGroup.Filters = append(fGroup.Filters, Filter{
		Value:        phoneNumber,
		PropertyName: "phone",
		Operator:     "EQ",
	})

	filterQuery.FilterGroups = append(filterQuery.FilterGroups, fGroup)

	marshal, err := json.Marshal(filterQuery)
	if err != nil {
		log.Fatalf("%s", "Error parsing json.")
	}
	log.Infof("%s", marshal)

	resp := ContactsSearchResponse{}

	err = c.Client.Request("POST", "/crm/v3/objects/contacts/search", &filterQuery, &resp)
	return resp, err
}
