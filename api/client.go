package api

import (
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"io/ioutil"
)

// Client structure
type Client struct {
	client *http.Client

	c3poInstance string
	c3poApplicationID string
	c3poUsername string
	c3poPassword string
	c3poAccessToken string

	debug bool
}

const (
        keystoneTarget        = "prod"
        keystoneAPIServer     = "https://api.keystone.disney.com"
        keystoneApplication   = "TWDC.ParksandResorts.c3po-prod"
        keystoneApplicationID = "4515ed23-5479-4cb0-a342-817b90e21241"
)

// NewClient - Creates a new client and returns an error if it fails
func NewClient(c3poUsername, c3poPassword, c3poAccessToken string, debug bool) (*Client, error) {
	if c3poUsername == "" || c3poPassword == "" {
		return nil, fmt.Errorf("username or password cannot be empty")
	}

	c := &Client{
		c3poInstance:      keystoneAPIServer,
		c3poApplicationID: keystoneApplicationID,
		c3poUsername:      c3poUsername,
		c3poPassword:      c3poPassword,
		c3poAccessToken:   c3poAccessToken,
		debug:             debug,
		client:            &http.Client{
			           Transport: &http.Transport{
				        TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			           },
		},
	}

	// Optionally, add more checks or initialization logic here

	return c, nil
}

func (c *Client) API(method string, apiHeader map[string][]string, apiEndpoint string, apiQueryString string, apiBody io.Reader) ([]byte, int, error) {

	if c.debug {
		fmt.Println("\n++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++")
		fmt.Println("++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++")
		fmt.Println("FUNC - API() Method : ", method)
		fmt.Println("FUNC - API() Header : ", apiHeader)
		fmt.Println("FUNC - API() Endpoint : ", apiEndpoint)
		fmt.Println("FUNC - API() Query : ", apiQueryString)
		fmt.Println("FUNC - API() Body : ", apiBody)
	}

	if method == "" {
		method = "GET"
	}
	apiMethod := strings.ToUpper(method)

	if apiHeader == nil {
		apiHeader = make(map[string][]string)
		apiHeader["Accept"] = []string{"application/json"}
		apiHeader["Content-Type"] = []string{"application/json"}
	}

	apiURL := c.c3poInstance + "/" + apiEndpoint

	// Initialize url.Values for query parameters
	apiQuery := url.Values{}
	if apiQueryString != "" {
		// Here you should parse the apiQueryString into url.Values
		// This depends on how your query string is structured.
		// For now, I'll assume it's in a standard format.
		// If it's not, you'll need to parse it according to your specific format.
		parsedQuery, err := url.ParseQuery(apiQueryString)
		if err != nil {
			// handle error
			log.Fatal("Query parsing error:", err)
		}
		apiQuery = parsedQuery
		// Append query parameters to the URL
		apiURL = apiURL + "?" + apiQuery.Encode()
	//} else {
		// If no query string is provided, set default values
		//apiQuery.Add("ApplicationId", c.c3poApplicationID)
		//apiQuery.Add("Directory", "vds")
		//apiQuery.Add("Username", c.c3poUsername)
		//apiQuery.Add("Password", c.c3poPassword)
	}


	if c.debug {
		first, last := extractFirstAndLast(c.c3poPassword)

		fmt.Println("----------------------------------------------------------------------------------------------------------")
		fmt.Println("API Method : ", apiMethod)
		fmt.Println("API URL : ", apiURL)
		fmt.Println("API Header : ", apiHeader)
		fmt.Println("API UserName/Password : [", c.c3poUsername, "]/[", first, "****", last, "]")
		fmt.Println("API BODY :", apiBody)
		fmt.Println("----------------------------------------------------------------------------------------------------------")
	}

	httpReq, httpReqErr := http.NewRequest(apiMethod, apiURL, apiBody)
	if httpReqErr != nil {
		log.Fatal("HTTP request creation error:", httpReqErr)
	}

	// Set headers from apiHeader to the request
	for key, values := range apiHeader {
		for _, value := range values {
			httpReq.Header.Add(key, value)
		}
	}

	//httpReq.SetBasicAuth(c.c3poUsername, c.c3poPassword)
	//if c.debug {
	//	fmt.Println(apiHeader)
	//	fmt.Println(httpReq.URL)
	//}

	if (c.debug) {
		fmt.Println(httpReq)
		fmt.Println("++++++++++++++++++++++")
	}

	httpResp, httpRespErr := c.client.Do(httpReq)
	if httpRespErr != nil {
		return nil, 0, httpRespErr
	} 

	defer httpResp.Body.Close()

	body, err := ioutil.ReadAll(httpResp.Body)
	if err != nil {
		return nil, httpResp.StatusCode, err
	}

	return body, httpResp.StatusCode, nil
}

func extractFirstAndLast(input string) (first, last string) {
    // Determine the length of the string
    length := len(input)

    // Set the number of characters to extract
    n := 4
    if length < n {
        n = length
    }

    // Extract the first and last characters
    first = input[:n]
    last = input[length-n:]

    return first, last
}
