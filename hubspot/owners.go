package hubspot

import "time"

// HSOwners client
type HSOwners struct {
	Client
}

type Owner struct {
	Id        string    `json:"id"`
	Email     string    `json:"email"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	UserId    int       `json:"userId"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	Archived  bool      `json:"archived"`
}

type OwnersResult struct {
	Results []Owner `json:"results"`
}

// Owners constructor (from Client)
func (c Client) Owners() HSOwners {
	return HSOwners{
		Client: c,
	}
}

func (c HSOwners) GetOwners() (OwnersResult, error) {
	resp := OwnersResult{}

	err := c.Client.Request("GET", "/crm/v3/owners", nil, &resp)
	return resp, err
}
