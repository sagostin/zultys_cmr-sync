package hubspot

import (
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
	FirstName        string    `json:"firstname"`
	LastName         string    `json:"lastname"`
	FullName         string    `json:"full_name"`
	LastModifiedDate time.Time `json:"lastmodifieddate"`
	Phone            string    `json:"phone"`
	Website          string    `json:"website"`
	SearchablePhone  string    `json:"hs_calculated_phone_number"`
}

type ContactResult struct {
	Id         string            `json:"id"`
	Properties ContactProperties `json:"properties"`
	CreatedAt  time.Time         `json:"createdAt"`
	UpdatedAt  time.Time         `json:"updatedAt"`
	Archived   bool              `json:"archived"`
}

type ContactsSearchResponse struct {
	Total   int             `json:"total"`
	Results []ContactResult `json:"results"`
	Paging  struct {
		Next struct {
			After string `json:"after"`
			Link  string `json:"link"`
		} `json:"next"`
	} `json:"paging"`
}

func (c Contacts) Get() (ContactsSearchResponse, error) {
	filterQuery := FilterQuery{
		FilterGroups: []FilterGroup{},
	}

	resp := ContactsSearchResponse{}

	err := c.Client.Request("POST", "/crm/v3/objects/contacts", &filterQuery, &resp)
	return resp, err
}

func (c Contacts) SearchByPhone(phoneNumber string) (ContactsSearchResponse, error) {
	filterQuery := FilterQuery{
		FilterGroups: []FilterGroup{},
		Properties:   []string{"hs_calculated_phone_number", "phone", "company", "full_name", "firstname", "lastname", "website"},
	}

	fGroup := FilterGroup{}
	fGroup.Filters = []Filter{}
	fGroup.Filters = append(fGroup.Filters, Filter{
		Value:        phoneNumber,
		PropertyName: "hs_calculated_phone_number",
		Operator:     "EQ",
	})

	filterQuery.FilterGroups = append(filterQuery.FilterGroups, fGroup)

	resp := ContactsSearchResponse{}

	err := c.Client.Request("POST", "/crm/v3/objects/contacts/search", &filterQuery, &resp)
	return resp, err
}
