package main

import (
	"encoding/json"
	"github.com/zultys_crm-sync/hubspot"
	"io/ioutil"
)

func SaveOwnersToFile(users []hubspot.Owner, filename string) error {
	data, err := json.Marshal(users)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filename, data, 0644)
}

func LoadOwnersFromFile(filename string) ([]hubspot.Owner, error) {
	var contacts []hubspot.Owner
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(data, &contacts)
	return contacts, err
}

// FindOwnerByEmail searches for an owner by email.
func FindOwnerByEmail(owners []hubspot.Owner, email string) *hubspot.Owner {
	for _, owner := range owners {
		if owner.Email == email {
			return &owner
		}
	}
	return nil
}

// FindOwnerByName searches for an owner by first and last name.
func FindOwnerByName(owners []hubspot.Owner, firstName, lastName string) *hubspot.Owner {
	for _, owner := range owners {
		if owner.FirstName == firstName && owner.LastName == lastName {
			return &owner
		}
	}
	return nil
}
