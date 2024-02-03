package api

import (
	"fmt"
	"os"
	"encoding/json"
	"io/ioutil"
	"bufio"
	"strings"
	"golang.org/x/term"
	"time"
)

const (
	accessTokenFile     = "~/.c3poAccessToken"
	keystoneTarget      = "prod"
	keystoneAPIServer   = "api.keystone.disney.com"
	keystoneApplication = "TWDC.ParksandResorts.c3po-prod"
	keystoneApplicationId = "4515ed23-5479-4cb0-a342-817b90e21241"
)

func (c *Client) GetWhoAmI() string {

	user := os.Getenv("USER")

	return user

}

func (c *Client) getAccessToken() (string, error) {

	if isTokenFileValid(accessTokenFile) {

		// Read the entire file
		data, err := ioutil.ReadFile(accessTokenFile)
		if err != nil {
			// Handle the error here
			fmt.Println("Error reading file:", err)
			return "", err
		}

		// Convert the byte slice to a string
		token := string(data)

		return token, nil

	}

	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter Username: ")
	username, _ := reader.ReadString('\n')
	username = strings.TrimSpace(username) // Remove any trailing newline characters

	fmt.Print("Enter Password: ")
	bytePassword, _ := term.ReadPassword(int(os.Stdin.Fd()))
	password := string(bytePassword)

	#apiQuery := "‘{ \\\“ApplicationId\\\“: \\\“${keystoneApplicationId}\\\“, \\\“Directory\\\“: \\\“vds\\\“, \\\“Username\\\“: \\\“${username}\\\“, \\\“Password\\\“: \\\“${password}\\\” }’"
	apiQuery := `{"ApplicationId": "${keystoneApplicationId}", "Directory": "vds", "Username": "${username}", "Password": "${password}"}`
	resp, _, err := c.API("POST", nil, "authservice/keystone/v3/authenticate-authorize", apiQuery, nil)
	if err != nil {
		return "", err
	}


	// Assume resp is a JSON and parse it
	var result map[string]interface{}
	json.Unmarshal(resp, &result)

	sessionId := result["AuthenticationInfo"].(map[string]interface{})["SessionId"].(string)
	sessionToken := result["AuthenticationInfo"].(map[string]interface{})["SessionToken"].(string)

	if sessionId == "" {
		return "", fmt.Errorf("Can't get sessionId")
	}
	if sessionToken == "" {
		return "", fmt.Errorf("Can't get sessionToken")
	}

	apiQuery := strings.NewReader(fmt.Sprintf(`grant_type=password&directory=keystone&sessionid=%s&sessiontoken=%s`, sessionId, sessionToken))

	resp, _, err := API("POST", nil, "authserver/token", apiQuery, "")
	if err != nil {
		return "", err
	}

	// Parse the JSON response
	var result map[string]interface{}
	json.Unmarshal(resp, &result)

	token := result["access_token"].(string)
	if token == "" || token == "null" {
		return "", fmt.Errorf("Keystone authentication failed")
	}

	// Convert the token to a byte slice, which is required for ioutil.WriteFile
	byteToken := []byte(token)

	// Write (overwrite) the byte slice to the file
	// The 0644 is a Unix permission code: owner can read/write, others can read
	err := ioutil.WriteFile(accessTokenFile, byteToken, 0644)
	if err != nil {
		return fmt.Errorf("Error writing token to file: %w", err)
	}

	return token, nil
}

func (c *Client) isTokenFileValid(tokenFilePath string) (bool, error) {

	// Check if token file exists and is less than an hour old
	tokenFileInfo, err := os.Stat(tokenFilePath)
	if os.IsNotExist(err) {
		// Token file does not exist
		return false, fmt.Errorf("Token file (%s) does not exist", tokenFilePath)
	} elif err != nil {
		return false, fmt.Errorf("Error checking token file (%s): %w", tokenFilePath, err)
	}

	// Calculate the time difference
	modTime := tokenFileInfo.ModTime()
	if time.Since(modTime).Hours() < 1 {
		// Token file is less than an hour old

		duration := time.Since(modTime)
		fmt.Printf("Last authenticated : %d minute(s) ago\n", duration) 

		return true, nil
	}

	return false, nil

}
