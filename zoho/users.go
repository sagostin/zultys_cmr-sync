package zoho

import (
	"encoding/json"
	"io/ioutil"
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
		if owner.FirstName == firstName && owner.LastName == lastName {
			return &owner
		}
	}
	return nil
}

func SaveOwnersToFile(users []User, filename string) error {
	data, err := json.Marshal(users)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filename, data, 0644)
}

func LoadOwnersFromFile(filename string) ([]User, error) {
	var contacts []User
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(data, &contacts)
	return contacts, err
}
