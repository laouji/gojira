package gojira

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

// ErrNoRepositories indicates no repos found
var ErrNoRepositories = errors.New("no repositories associated with this issue")

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
	res, err := client.getRequest(path, http.StatusOK)
	if err != nil {
		return issues, fmt.Errorf("IssuesByCustomFilter failed request: %w", err)
	}
	defer res.Body.Close()

	var result struct {
		Issues []*Issue `json:"issues"`
	}

	bodyBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return issues, fmt.Errorf("IssuesByCustomFilter failed to read response body: %w", err)
	}

	if err := json.Unmarshal(bodyBytes, &result); err != nil {
		return issues, fmt.Errorf("IssuesByCustomFilter failed unmarshal response body: %w", err)
	}
	return result.Issues, nil
}

// RepositoryType fetches the repository type associated with a particular issue
func (client *Client) RepositoryType(issueID string) (repositoryType string, err error) {
	path := fmt.Sprintf("/rest/dev-status/latest/issue/summary?issueId=%s", issueID)
	res, err := client.getRequest(path, http.StatusOK)
	if err != nil {
		return "", fmt.Errorf("RepositoryType failed request: %w", err)
	}
	defer res.Body.Close()

	var result DevStatus
	bodyBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("RepositoryType failed to read response body: %w", err)
	}

	if err := json.Unmarshal(bodyBytes, &result); err != nil {
		return "", fmt.Errorf("RepositoryType failed unmarshal response body: %w", err)
	}

	if len(result.Errors) > 0 {
		return "", fmt.Errorf("RepositoryType found unexpected errors: %s", result.Errors)
	}

	instances := result.Summary.Branch.ByInstanceType
	if len(instances) == 0 {
		instances = result.Summary.Repository.ByInstanceType
		if len(instances) == 0 {
			return "", ErrNoRepositories
		}
	}
	if len(instances) > 1 {
		return "", fmt.Errorf("RepositoryType expected 1 repository type but found %d", len(instances))
	}

	for key := range instances {
		repositoryType = key
	}
	return repositoryType, nil
}

// Branches fetches info about the branches associated with this issuue
func (client *Client) Branches(issueID, repositoryType string) (branches []*Branch, err error) {
	path := fmt.Sprintf("/rest/dev-status/latest/issue/detail?issueId=%s&applicationType=%s&dataType=branch", issueID, repositoryType)
	res, err := client.getRequest(path, http.StatusOK)
	if err != nil {
		return []*Branch{}, fmt.Errorf("Branches failed request: %w", err)
	}
	defer res.Body.Close()

	var result DevStatus
	bodyBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return []*Branch{}, fmt.Errorf("Branches failed to read response body %w", err)
	}

	if err := json.Unmarshal(bodyBytes, &result); err != nil {
		return []*Branch{}, fmt.Errorf("Branches failed unmarshal response body: %w", err)
	}

	if len(result.Errors) > 0 {
		return []*Branch{}, fmt.Errorf("Branches found unexpected errors: %s", result.Errors)
	}
	branches = make([]*Branch, 0)
	for _, detail := range result.Details {
		branches = append(branches, detail.Branches...)
	}
	return branches, nil
}

func (client *Client) getRequest(path string, expectedStatus int) (res *http.Response, err error) {
	res, err = client.httpClient.Get(client.baseURL.String() + path)
	if err != nil {
		return nil, fmt.Errorf("failed to make request to %s: %w", path, err)
	}

	if res.StatusCode != expectedStatus {
		return nil, fmt.Errorf("request to %s responded with status %d", path, res.StatusCode)
	}
	return res, nil
}
