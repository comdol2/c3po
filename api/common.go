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

type Group struct {
    Id                     string       `json:"Id"`
    Name                   string       `json:"Name"`
    DynamicAssignmentId    interface{}  `json:"DynamicAssignmentId"`
    LastUpdate             string       `json:"LastUpdate"`
}

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

type GroupAttributes struct {
    AttributeId            string       `json:"AttributeId"`
    AttributeName          string       `json:"AttributeName"`
    AttributeValue         string       `json:"AttributeValue"`
    GroupId                string       `json:"GroupId"`
    GroupRoleId            string       `json:"GroupRoleId"`
    Id                     string       `json:"Id"`
    LastUpdate             string       `json:"LastUpdate"`
    RoleId                 string       `json:"RoleId"`
    UsageType              int          `json:"UsageType"`
}

type User struct {
    CommonName             string       `json:"CommonName"`
    Ecrid                  string       `json:"Ecrid"`
    Email                  string       `json:"Email"`
    FirstName              string       `json:"FirstName"`
    Id                     string       `json:"Id"`
    IdAtSourceSystem       string       `json:"IdAtSourceSystem"`
    IsActive               bool         `json:"IsActive"`
    LastName               string       `json:"LastName"`
    SourceSystemId         string       `json:"SourceSystemId"`
    SourceSystemName       string       `json:"SourceSystemName"`
}

type FunctionalAbilities struct {
    ApplicationId                       string      `json:"ApplicationId"`
    DataClassification                  int         `json:"DataClassification"`
    Description                         string      `json:"Description"`
    FunctionalAbilityEntityAccess       interface{} `json:"FunctionalAbilityEntityAccess"`
    Id                                  string      `json:"Id"`
    LastUpdate                          string      `json:"LastUpdate"`
    Name                                string      `json:"Name"`
    SodRole                             string      `json:"SodRole"`

}

const (
	accessTokenFile     = "~/.c3poAccessToken"
)

// GetAbsolutePath takes a path string and returns its absolute path.
func (c *Client) GetAbsolutePath(path string) (string, error) {
    // Expand the '~' if used
    if len(path) > 0 && path[:1] == "~" {
        usr, err := user.Current()
        if err != nil {
            return "", fmt.Errorf("failed to get current user: %w", err)
        }
        path = filepath.Join(usr.HomeDir, path[1:])
    }

    // Get absolute path
    absPath, err := filepath.Abs(path)
    if err != nil {
        return "", fmt.Errorf("failed to get absolute path: %w", err)
    }

    return absPath, nil
}

func (c *Client) removeC3POPrefixes(input string) string {
	// Define the prefixes to remove
	prefixes := []string{"C3PO -", "C3PO-", "C3PO", "C3PO\\s+-"}

	// Construct the regular expression pattern
	pattern := "(" + strings.Join(prefixes, "|") + ")\\s*$"

	// Compile the regular expression
	regexpPattern := regexp.MustCompile(pattern)

	// Replace matching prefixes with an empty string
	result := regexpPattern.ReplaceAllString(input, "")

	return strings.TrimSpace(result)
}

func (c *Client) reduceSpaces(input string) string {
	// Define the regular expression pattern to match two or more spaces
	regex := regexp.MustCompile(`\s{2,}`)

	// Replace multiple spaces with a single space
	result := regex.ReplaceAllString(input, " ")

	return result
}

// Define custom sorting function by role ID
func (c *Client) SortByRoleID(roles []Role) {
        sort.Slice(roles, func(i, j int) bool {
                return roles[i].Name < roles[j].Name
        })
}

// Print roles neatly
func (c *Client) PrintRoles(roles []Role) {
	fmt.Println("===========================================\n")
	totalRoles := len(roles)
	for i, role := range roles {
		index := i + 1
		if ( i > 0 ) {
			fmt.Println("-------------------------------------------")
		}
		fmt.Printf("Role %d/%d: %s\n", index, totalRoles, role.Name)
		fmt.Println("\tApplicationId:", role.ApplicationId)
		fmt.Println("\tDescription:", role.Description)
		fmt.Println("\tConditionalExpression:", role.ConditionalExpression)
		// Add other fields here
		fmt.Println()
	}
}

// Print groups neatly
func (c *Client) PrintGroups(groups []Group) {
        for _, group := range groups {
                //fmt.Println("Group ID:", group.Id)
                fmt.Println("Name: " + group.Name)
                //fmt.Println("Description:", group.Description)
                //fmt.Println("Last Update:", group.LastUpdate)
                // Print other fields as needed
                //fmt.Println("----------------------------------")
        }
}

// Print groups neatly
func (c *Client) PrintFunctionalAbilities(functionalabilities []FunctionalAbilities) {
        for _, functionalability := range functionalabilities {
                fmt.Println("Name: " + functionalability.Name)
        }
}

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
