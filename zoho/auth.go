package zoho

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	_ "strings"
	"time"
)

type Endpoints struct {
	AccountAuth string `json:"account_auth"`
	CrmApi      string `json:"crm_api"`
}

type AccessGrant struct {
	Scope        string `json:"scope,omitempty"`
	ExpiryTime   int64  `json:"expires_in,omitempty"`
	ClientId     string `json:"client_id,omitempty"`
	ClientSecret string `json:"client_secret,omitempty"`
	AccessToken  string `json:"access_token,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
	GrantCode    string `json:"code,omitempty"`
	GrantType    string `json:"grant_type,omitempty"`
	TokenType    string `json:"token_type,omitempty"`
}

type Client struct {
	Auth      AccessGrant `json:"auth"`
	Endpoints Endpoints   `json:"endpoints"`
}

func (c *Client) Authenticate(clientId, clientSecret, grantCode string) error {
	// Construct the form data from the AccessGrant struct
	data := url.Values{}
	data.Set("client_id", clientId)
	data.Set("client_secret", clientSecret)
	data.Set("code", grantCode)
	data.Set("grant_type", "authorization_code")
	// If scope is required as a space-separated list
	//data.Set("scope", strings.Join(c.Auth.Scope, " "))

	// Construct the request
	req, err := http.NewRequest("POST", "https://"+c.Endpoints.AccountAuth+"/oauth/v2/token", strings.NewReader(data.Encode()))
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %v", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Error(err)
		}
	}(resp.Body)

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %v", err)
	}
	log.Println(string(body))

	// TODO: Extract the token from the response and start the refresher thread
	// For now, just log the body
	log.Printf("Authenticate response: %s", body)

	//var auth AccessGrant

	err = json.Unmarshal(body, &c.Auth)
	if err != nil {
		log.Error(err)
		return err
	}

	//c.Auth = auth

	log.Println("Authentication successful, token acquired and refresher started.")

	return nil
}

func (c *Client) StartTokenRefresher() {
	// Calculate the duration to wait before refreshing the token
	// For example, if the token expires in 3600 seconds, you might want to refresh after 3500 seconds

	waitDuration := c.Auth.ExpiryTime - 100

	// Start a new goroutine to handle the refreshing
	go func() {
		for {
			time.Sleep(time.Duration(waitDuration) * time.Second)
			err := c.refreshAccessToken()
			if err != nil {
				log.Printf("Error refreshing token: %v", err)
				// Handle the error, possibly by trying to refresh the token again after a short delay
			} else {
				log.Info("Token refreshed successfully")
				// Token refreshed successfully, calculate the next wait duration based on the new expiry time
				// You would update waitDuration here based on the new expiry time received in RefreshToken
			}
		}
	}()
}

// TokenResponse represents the JSON structure of the response from Zoho OAuth server
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int64  `json:"expires_in"`
	ApiDomain    string `json:"api_domain"`
	TokenType    string `json:"token_type"`
	RefreshToken string `json:"refresh_token"`
}

func (c *Client) refreshAccessToken() error {
	// Prepare the query parameters
	data := url.Values{}
	data.Set("client_id", c.Auth.ClientId)
	data.Set("client_secret", c.Auth.ClientSecret)
	data.Set("refresh_token", c.Auth.RefreshToken)
	data.Set("grant_type", "refresh_token")
	// If scope is required as a space-separated list
	//data.Set("scope", strings.Join(c.Auth.Scope, " "))

	// Construct the request
	req, err := http.NewRequest("POST", "https://"+c.Endpoints.AccountAuth+"/oauth/v2/token", strings.NewReader(data.Encode()))
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %v", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Error(err)
		}
	}(resp.Body)

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %v", err)
	}
	log.Println(string(body))

	// TODO: Extract the token from the response and start the refresher thread
	// For now, just log the body
	log.Printf("Authenticate response: %s", body)

	err = json.Unmarshal(body, &c.Auth)
	if err != nil {
		log.Error(err)
		return err
	}

	// Log the successful refresh
	fmt.Printf("Successfully refreshed token. New access token: %s\n", c.Auth.AccessToken)

	return nil
}

// depending on the dc the url will be different, so we will need to have all the info, split by colons?
// eg. https://accounts.zohocloud.ca/oauth/v2/token:https://crm.zohocloud.ca/crm/v6:CLIENT_ID:CLIENT_SECRET:GRANT
