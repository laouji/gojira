package gojira_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/laouji/gojira"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIssuesByFilter(t *testing.T) {
	server := httptest.NewServer(testHandler())
	defer server.Close()
	client, err := gojira.NewClient(server.URL, "user", "password", 2*time.Second)
	require.NoError(t, err)

	issues, err := client.IssuesByCustomFilter("10000", "COR-3425")
	require.NoError(t, err)
	require.Len(t, issues, 1)
	assert.Equal(t, "123456", issues[0].ID)
	assert.Equal(t, "PRO-1234", issues[0].Key)
}

func testHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"issues":[{"id":"123456","key":"PRO-1234"}]}`))
	})
}
