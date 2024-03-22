package api

import (
	"fmt"
	"strings"
	"os/user"
	"path/filepath"
	"regexp"
	"sort"
)

type Group struct {
    Id                     string       `json:"Id"`
    Name                   string       `json:"Name"`
    DynamicAssignmentId    interface{}  `json:"DynamicAssignmentId"`
    LastUpdate             string       `json:"LastUpdate"`
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


















