package main

import (
	"encoding/json"
	"fmt"
	zultys "github.com/sagostin/zultys-go/lib"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"reflect"
	"sort"
	"strings"
)

type ZultysUser struct {
	CallerId  string `json:"callerId"`
	CellPhone string `json:"cellPhone"`
	Did       string `json:"did"`
	Email     string `json:"email"`
	Extension string `json:"extension"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

func getZultysUsers(config Config) ([]ZultysUser, error) {
	loginRequest := zultys.LoginRequest{
		UserLogin: config.MxUsername,
		Password:  config.MxPassword,
	}

	zSess, _, err := loginRequest.Login(config.MxAddr)
	if err != nil {
		return nil, err
	}

	getParams := map[string]interface{}{
		"page":  "1",
		"limit": "100",
	}

	tt, err := zSess.SendCommand(zultys.MethodGet, "adm_get_users", getParams, nil)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	response, err := zultys.HandleGetUsersResponse(tt)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	var users []ZultysUser

	for _, rUser := range response.Users {
		user := ZultysUser{}

		marshal, err := json.Marshal(rUser)
		if err != nil {
			log.Error(err)
			return nil, err
		}

		err = json.Unmarshal(marshal, &user)
		if err != nil {
			log.Error(err)
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}

type ZultysUsers []ZultysUser

func SaveUsersToFile(users ZultysUsers, filename string) error {
	data, err := json.Marshal(users)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filename, data, 0644)
}

func LoadUsersFromFile(filename string) (ZultysUsers, error) {
	var contacts ZultysUsers
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return contacts, err
	}
	err = json.Unmarshal(data, &contacts)
	return contacts, err
}

func UpdateContactsFile(newContacts ZultysUsers, filename string) error {
	// Load existing contacts
	existingContacts, err := LoadUsersFromFile(filename)
	if err != nil {
		return err
	}

	// Compare existing with new data here and update as necessary...
	// This example assumes you replace the entire list, but you might
	// want to update individual records based on your logic.
	if !reflect.DeepEqual(existingContacts, newContacts) {
		fmt.Println("Updating contacts...")
		return SaveUsersToFile(newContacts, filename)
	}

	fmt.Println("No updates necessary.")
	return nil
}

// SortEntriesByExtension sorts the entries by their extension.
func SortEntriesByExtension(entries []ZultysUser) {
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Extension < entries[j].Extension
	})
}

// FindEntryByExtension searches for an entry by extension.
func FindEntryByExtension(entries []ZultysUser, extension string) *ZultysUser {
	for _, entry := range entries {
		if entry.Extension == extension {
			return &entry
		}
	}
	return nil
}

// FindEntryByName searches for an entry by first and last name.
func FindEntryByName(entries []ZultysUser, firstName, lastName string) *ZultysUser {
	for _, entry := range entries {
		if strings.EqualFold(entry.FirstName, firstName) && strings.EqualFold(entry.LastName, lastName) {
			return &entry
		}
	}
	return nil
}
