package gojira

// DevStatus is a json representation of http Responses from the dev-status plugin
type DevStatus struct {
	Errors  []string          `json:"errors"`
	Summary devStatusSummary  `json:"summary,omitempty"`
	Details []devStatusDetail `json:"detail,omitempty"`
}

// Branch is a json representation of branch info given by the dev-status plugin
type Branch struct {
	Name       string `json:"name"`
	URL        string `json:"url"`
	Repository struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"repository"`
}

type devStatusSummary struct {
	Repository struct {
		ByInstanceType map[string]map[string]interface{} `json:"byInstanceType"`
	} `json:"repository"`
}

type devStatusDetail struct {
	Repositories []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"repositories"`
	Branches []*Branch `json:"branches"`
}
