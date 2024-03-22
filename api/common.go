package api

import (
	"bytes"
	"fmt"
	"os"
	"encoding/json"
	"time"
	"net/http"
	"net/url"
	"strings"
	"regexp"
)

type Role struct {
        ApplicationId          string      `json:"ApplicationId"`
        Name                   string      `json:"Name"`
        Description            string      `json:"Description"`
        ConditionalExpression  interface{} `json:"ConditionalExpression"`
        DynamicAssignmentId    interface{} `json:"DynamicAssignmentId"`
        RoleFunctionalAbilities interface{} `json:"RoleFunctionalAbilities"`
        Id                     string      `json:"Id"`
        LastUpdate             string      `json:"LastUpdate"`
}

const (
	accessTokenFile     = "~/.c3poAccessToken"
)

func (c *Client) GetAccessToken() (string, error) {
	apiBodyData := map[string]string{
		"ApplicationId": c.c3poApplicationID,
		"Directory":     "vds",
		"Username":      c.c3poUsername,
		"Password":      c.c3poPassword,
	}

	// Convert the map to a JSON string for the body
	jsonBody, err := json.Marshal(apiBodyData)
	if err != nil {
		return "", err
	}

	respBytes, statusCode, err := c.API("POST", nil, "authservice/keystone/v3/authenticate-authorize", "", bytes.NewReader(jsonBody))
	if err != nil {
		return "", err
	}

	if (c.debug) {
		// Debug: Print the response as string
		fmt.Println("authenticate-authorize => Response as string:", string(respBytes))
	}

	if statusCode != http.StatusOK {
		return "", fmt.Errorf("API request failed with status code: %d", statusCode)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(respBytes, &result); err != nil {
		return "", err
	}


	if (c.debug) {
		fmt.Println(result)
	}

	sessionId := result["AuthenticationInfo"].(map[string]interface{})["SessionId"].(string)
	sessionToken := result["AuthenticationInfo"].(map[string]interface{})["SessionToken"].(string)

	if sessionId == "" {
		return "", fmt.Errorf("Can't get sessionId")
	}
	if sessionToken == "" {
		return "", fmt.Errorf("Can't get sessionToken")
	}

	if (c.debug) {
		fmt.Println(">>> sessionId: ", sessionId)
		fmt.Println(">>> sessionToken: ", sessionToken)
	}

	// Prepare URL-encoded form data
	formData := url.Values{}
	formData.Set("grant_type", "password")
	formData.Set("directory", "keystone")
	formData.Set("sessionid", sessionId)
	formData.Set("sessiontoken", sessionToken)

	apiQuery := strings.NewReader(formData.Encode()) // Convert form data to io.Reader

	respBytes, statusCode, err = c.API("POST", nil, "authserver/token", "", apiQuery)
	if err != nil {
		return "", err
	}

	if (c.debug) {
		// Debug: Print the response as string
		fmt.Println("authserver/token => Response as string:", string(respBytes))
	}

	if statusCode != http.StatusOK {
		return "", fmt.Errorf("API request failed with status code: %d", statusCode)
	}

	if err := json.Unmarshal(respBytes, &result); err != nil {
		return "", err
	}

	token, ok := result["access_token"].(string)
	if !ok || token == "" {
		return "", fmt.Errorf("Keystone authentication failed")
	}

	c.c3poAccessToken = token

	absAccessTokenFile, err1 := c.GetAbsolutePath(accessTokenFile)
        if err1 != nil {
		return "", fmt.Errorf("Error: %w", err1)
        }

	byteToken := []byte(token)
	err = os.WriteFile(absAccessTokenFile, byteToken, 0600)
	if err != nil {
		return "", fmt.Errorf("error writing token to file: %w", err)
	}

	return token, nil

}

func (c *Client) IsTokenFileValid(tokenFilePath string) (bool, error) {

	// Check if token file exists and is less than an hour old
	tokenFileInfo, err := os.Stat(tokenFilePath)
	if os.IsNotExist(err) {
		// Token file does not exist
		return false, fmt.Errorf("Token file (%s) does not exist", tokenFilePath)
	} else if err != nil {
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

func (c *Client) GetGroup(rolename string, exactmatched bool) ([]map[string]interface{}, error) {
	// FilteredItems to store the filtered roles
	var filteredItems []map[string]interface{}

	requestedGroupEscaped := url.QueryEscape("C3PO - " + rolename)
	if c.debug {
		fmt.Println("Requested Group Escaped:", requestedGroupEscaped)
	}

	KeyStoneAPIPath := "adminservice/keystone/v1/group?groupName=" + requestedGroupEscaped
	if c.debug {
		fmt.Println("KeyStone API Path:", KeyStoneAPIPath)
	}

	REQ_METHOD := "GET"
	apiHeader := make(map[string][]string)
	apiHeader["Authorization"] = []string{"Bearer " + c.c3poAccessToken}
	apiHeader["Accept"] = []string{"application/json"}

	respBytes, statusCode, err := c.API(REQ_METHOD, apiHeader, KeyStoneAPIPath, "", nil)
	if err != nil {
		return nil, err
	}

	if c.debug {
		fmt.Println("Get Role => Response as string:", string(respBytes))
	}

	if statusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status code: %d", statusCode)
	}

	var roles []map[string]interface{}
	if err := json.Unmarshal(respBytes, &roles); err != nil {
		return nil, err
	}

	if exactmatched {
		requestedGroupLower := strings.ToLower("C3PO - " + rolename)
		re := regexp.MustCompile(`[()]`)
		requestedGroupLower = re.ReplaceAllStringFunc(requestedGroupLower, func(s string) string {
			return "\\" + s
		})
		if c.debug {
			fmt.Println("Requested Group Lower:", requestedGroupLower)
		}

		pattern := fmt.Sprintf("^%s$", regexp.QuoteMeta(requestedGroupLower))

		for _, item := range roles {
			name, ok := item["Name"].(string)
			if !ok {
				continue
			}
			if matched, _ := regexp.MatchString(pattern, strings.ToLower(name)); matched {
				filteredItems = append(filteredItems, item)
			}
		}
	} else {
		roleNameLower := strings.ToLower(rolename)
		for _, item := range roles {
			name, ok := item["Name"].(string)
			if !ok {
				continue
			}
			if strings.Contains(strings.ToLower(name), roleNameLower) {
				filteredItems = append(filteredItems, item)
				break
			}
		}
	}

	return filteredItems, nil
}

func (c *Client) GetRole(rolename string, exactmatched bool) ([]map[string]interface{}, error) {
	// FilteredItems to store the filtered roles
	var filteredItems []map[string]interface{}

	KeyStoneAPIPath := "adminservice/keystone/v1/application/" + c.c3poApplicationID + "/role"
	if c.debug {
		fmt.Println("KeyStone API Path:", KeyStoneAPIPath)
	}

	REQ_METHOD := "GET"
	apiHeader := make(map[string][]string)
	apiHeader["Authorization"] = []string{"Bearer " + c.c3poAccessToken}
	apiHeader["Accept"] = []string{"application/json"}

	respBytes, statusCode, err := c.API(REQ_METHOD, apiHeader, KeyStoneAPIPath, "", nil)
	if err != nil {
		return nil, err
	}

	if c.debug {
		fmt.Println("Get Role => Response as string:", string(respBytes))
	}

	if statusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status code: %d", statusCode)
	}

	var roles []map[string]interface{}
	if err := json.Unmarshal(respBytes, &roles); err != nil {
		return nil, err
	}

	if c.debug {
		fmt.Println(roles)
	}

	roleNameLower := strings.ToLower(rolename)
	if exactmatched {
		roleNameLower := strings.ToLower(rolename)
		re := regexp.MustCompile(`[()]`)
		roleNameLower = re.ReplaceAllStringFunc(roleNameLower, func(s string) string {
			return "\\" + s
		})
		if c.debug {
			fmt.Println("Requested Role Lower:", roleNameLower)
		}

		pattern := fmt.Sprintf("^%s$", regexp.QuoteMeta(roleNameLower))

		for _, item := range roles {
			name, ok := item["Name"].(string)
			if !ok {
				continue
			}
			if matched, _ := regexp.MatchString(pattern, strings.ToLower(name)); matched {
				filteredItems = append(filteredItems, item)
			}
		}
	} else {
		for _, item := range roles {
			name, ok := item["Name"].(string)
			if !ok {
				continue
			}
			if strings.Contains(strings.ToLower(name), roleNameLower) {
				filteredItems = append(filteredItems, item)
				break
			}
		}
	}

	return filteredItems, nil
}
