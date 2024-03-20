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

type User struct {
	Country              string    `json:"country"`
	Language             string    `json:"language"`
	Id                   string    `json:"id"`
	CountryLocale        string    `json:"country_locale"`
	DecimalSeparator     string    `json:"decimal_separator"`
	CreatedTime          time.Time `json:"created_time"`
	TimeFormat           string    `json:"time_format"`
	Offset               int       `json:"offset"`
	Zuid                 string    `json:"zuid"`
	FullName             string    `json:"full_name"`
	Phone                string    `json:"phone"`
	SortOrderPreferenceS string    `json:"sort_order_preference__s"`
	Status               string    `json:"status"`
	Locale               string    `json:"locale"`
	PersonalAccount      bool      `json:"personal_account,omitempty"`
	Isonline             bool      `json:"Isonline"`
	DefaultTabGroup      string    `json:"default_tab_group,omitempty"`
	FirstName            string    `json:"first_name"`
	Email                string    `json:"email"`
	ModifiedTime         time.Time `json:"Modified_Time"`
	Mobile               string    `json:"mobile"`
	LastName             string    `json:"last_name"`
	TimeZone             string    `json:"time_zone"`
	NumberSeparator      string    `json:"number_separator"`
	Confirm              bool      `json:"confirm"`
	DateFormat           string    `json:"date_format"`
	Category             string    `json:"category"`
}

type UsersResponse struct {
	Users []User `json:"users"`
	Info  struct {
		PerPage     int  `json:"per_page"`
		Count       int  `json:"count"`
		Page        int  `json:"page"`
		MoreRecords bool `json:"more_records"`
	} `json:"info"`
}

func (c *Client) FetchUsers() (UsersResponse, error) {
	var usersResponse UsersResponse

	// Make sure we have a valid access token
	if c.Auth.AccessToken == "" {
		if err := c.refreshAccessToken(); err != nil {
			return usersResponse, fmt.Errorf("failed to refresh access token: %v", err)
		}
	}

	// Construct the GET request to fetch users
	req, err := http.NewRequest("GET", "https://"+c.Endpoints.CrmApi+"/crm/v6/users", nil)
	if err != nil {
		return usersResponse, fmt.Errorf("failed to create request: %v", err)
	}

	// Set the Authorization header with the Bearer token
	req.Header.Add("Authorization", "Bearer "+c.Auth.AccessToken)

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return usersResponse, fmt.Errorf("failed to send request: %v", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Error(err)
		}
	}(resp.Body)

	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return usersResponse, fmt.Errorf("received non-OK response status: %s, body: %s", resp.Status, body)
	}

	// Decode the response body into the UsersResponse struct
	if err := json.NewDecoder(resp.Body).Decode(&usersResponse); err != nil {
		return usersResponse, fmt.Errorf("error decoding users response: %v", err)
	}

	return usersResponse, nil
}

// FindUserByEmail searches for an owner by email.
func FindUserByEmail(owners []User, email string) *User {
	for _, owner := range owners {
		if owner.Email == email {
			return &owner
		}
	}
	return nil
}

// FindUserByName searches for an owner by first and last name.
func FindUserByName(owners []User, firstName, lastName string) *User {
	for _, owner := range owners {
		if owner.FullName == firstName+" "+lastName {
			return &owner
		}
	}
	return nil
}

func SaveUsersToFile(users []User, filename string) error {
	data, err := json.Marshal(users)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filename, data, 0644)
}

func LoadUsersFromFile(filename string) ([]User, error) {
	var contacts []User
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(data, &contacts)
	return contacts, err
}
