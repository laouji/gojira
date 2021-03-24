package gojira

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

// Client is an API Client for making request to Jira
//
// it currently supports username:password style authentication
// which of course is not a very secure authentication method
// but seeing that jira doesn't support API tokens out of the box this'll have to do for now
type Client struct {
	baseURL    *url.URL
	httpClient *http.Client
}

// NewClient instantiates a jira.Client
func NewClient(
	baseURL string,
	userName string,
	password string,
	clientTimeout time.Duration,
) (*Client, error) {
	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse base url: %w", err)
	}
	u.User = url.UserPassword(userName, password)

	return &Client{
		baseURL:    u,
		httpClient: &http.Client{Timeout: clientTimeout},
	}, nil
}

// IssuesByCustomFilter fetches a list of issues that match the given custom filter
//
// TODO: might need paging
func (client *Client) IssuesByCustomFilter(filterName, filterValue string) (issues []*Issue, err error) {
	path := fmt.Sprintf("/rest/api/2/search?jql=cf[%s]=%s", filterName, filterValue)
	res, err := client.httpClient.Get(client.baseURL.String() + path)
	if err != nil {
		return issues, fmt.Errorf("failed to make search request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return issues, fmt.Errorf("search request responded with status %d", res.StatusCode)
	}

	var result struct {
		Issues []*Issue `json:"issues"`
	}

	bodyBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return issues, fmt.Errorf("failed to read response body: %w", err)
	}

	if err := json.Unmarshal(bodyBytes, &result); err != nil {
		return issues, fmt.Errorf("failed unmarshal response body: %w", err)
	}
	return result.Issues, nil
}
