package api

import (
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
)

// Client structure
type Client struct {
	client *http.Client

	c3poInstance string
	c3poUsername string
	c3poPassword string

	debug bool
}

// NewClient - Creates a new client
func NewClient(c3poUsername, c3poPassword string, debug bool) *Client {

	var c *Client

	if c3poUsername != "" && c3poPassword != "" {

		c = &Client{}

		c.c3poInstance = "https://ui.keystone.disney.com"
		c.c3poUsername = c3poUsername
		c.c3poPassword = c3poPassword

		transport := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		c.client = &http.Client{Transport: transport}
		c.debug = debug

	}

	return c

}

//method = "POST"
//apiHeader = nil
//apiEndpoint = "authservice/keystone/v3/authenticate-authorize"
//apiQuery = "'{ \\\"ApplicationId\\\": \\\"${KeyStone_ApplicationId}\\\", \\\"Directory\\\": \\\"vds\\\", \\\"Username\\\": \\\"${whoami}\\\", \\\"Password\\\": \\\"${password}\\\" }’"
//apiBody = ""

func (c *Client) API(method string, apiHeader map[string][]string, apiEndpoint string, apiQuery string, apiBody io.Reader) (apiresp interface{}, apirespcode int, apierr error) {

	if c.debug {
		fmt.Println("\n++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++")
		fmt.Println("FUNC - API Method : ", method)
		fmt.Println("FUNC - API Header : ", apiHeader)
		fmt.Println("FUNC - API Endpoint : ", apiEndpoint)
		fmt.Println("FUNC - API Query : ", apiQuery)
		fmt.Println("FUNC - API Body : ", apiBody)
	}

	if method == "" {
		method = "GET"
	}
	apiMethod := strings.ToUpper(method)

	if apiHeader == nil {
		apiHeader = url.Values{}
	}
	apiHeader["Accept"] = []string{"application/json"}
	apiHeader["Content-Type"] = []string{"application/json"}

	apiURL := c.c3poInstance + "/" + apiEndpoint
	if apiQuery != nil {
		apiURL = apiURL + " --data-raw " + apiQuery.Encode()
	}

	if c.debug {
		fmt.Println("----------------------------------------------------------------------------------------------------------")
		fmt.Println("API Method : ", apiMethod)
		fmt.Println("API URL : ", apiURL)
		fmt.Println("API Header : ", apiHeader)
		fmt.Println("API UserName/Password : [", c.c3poUsername, "]/[", c.c3poPassword, "]")
		fmt.Println("API BODY :", apiBody)
		fmt.Println("----------------------------------------------------------------------------------------------------------")
	}

	httpReq, httpReqErr := http.NewRequest(apiMethod, apiURL, apiBody)
	if httpReqErr != nil {
		log.Fatal("HTTP request creation error:", httpReqErr)
	}
	if c.debug {
		fmt.Println("Request URL: " + apiURL)
	}

	httpReq.SetBasicAuth(c.c3poUsername, c.c3poPassword)
	if c.debug {
		fmt.Println(apiHeader)
		fmt.Println(httpReq.URL)
	}

	statusCode := 0
	httpResp, httpRespErr := c.client.Do(httpReq)
	if httpRespErr != nil {
		fmt.Println("HTTP request error:", httpRespErr)
		fmt.Println("++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++\n")
	} else {
		defer httpResp.Body.Close()

		// Access the HTTP status code
		statusCode = httpResp.StatusCode
		fmt.Println("statusCode: ", statusCode)
		fmt.Println("++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++\n")

		if statusCode == 200 || statusCode == 201 {
			return httpResp, statusCode, nil
		}
	}
	return nil, statusCode, httpRespErr
}
