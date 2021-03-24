package gojira

// Issue is a json representation of a Jira Issue
type Issue struct {
	ID  string `json:"id"`
	Key string `json:"key"`
}
