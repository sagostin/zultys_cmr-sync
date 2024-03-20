package hubspot

import (
	"encoding/json"
	"io/ioutil"
)

func SaveOwnersToFile(users []Owner, filename string) error {
	data, err := json.Marshal(users)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filename, data, 0644)
}

func LoadOwnersFromFile(filename string) ([]Owner, error) {
	var contacts []Owner
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(data, &contacts)
	return contacts, err
}

// FindOwnerByEmail searches for an owner by email.
func FindOwnerByEmail(owners []Owner, email string) *Owner {
	for _, owner := range owners {
		if owner.Email == email {
			return &owner
		}
	}
	return nil
}

// FindOwnerByName searches for an owner by first and last name.
func FindOwnerByName(owners []Owner, firstName, lastName string) *Owner {
	for _, owner := range owners {
		if owner.FirstName == firstName && owner.LastName == lastName {
			return &owner
		}
	}
	return nil
}
