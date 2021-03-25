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
	resBody := []byte(`{"issues":[{"id":"123456","key":"PRO-1234"}]}`)
	server := httptest.NewServer(testHandler(resBody))
	defer server.Close()
	client, err := gojira.NewClient(server.URL, "user", "password", 2*time.Second)
	require.NoError(t, err)

	issues, err := client.IssuesByCustomFilter("10000", "PRO-3425")
	require.NoError(t, err)
	require.Len(t, issues, 1)
	assert.Equal(t, "123456", issues[0].ID)
	assert.Equal(t, "PRO-1234", issues[0].Key)
}

func TestRepositoryType_FromRepository(t *testing.T) {
	resBody := []byte(`{"errors":[],"summary":{"repository":{"byInstanceType":{"githube":{"count":15,"name":"GitHub Enterprise"}}}}}`)
	server := httptest.NewServer(testHandler(resBody))
	defer server.Close()
	client, err := gojira.NewClient(server.URL, "user", "password", 2*time.Second)
	require.NoError(t, err)

	rType, err := client.RepositoryType("112233")
	require.NoError(t, err)
	assert.Equal(t, "githube", rType)
}

func TestRepositoryType_FromBranch(t *testing.T) {
	resBody := []byte(`{"errors":[],"summary":{"branch":{"byInstanceType":{"githube":{"count":2,"name":"GitHub Enterprise"}}}}}`)
	server := httptest.NewServer(testHandler(resBody))
	defer server.Close()
	client, err := gojira.NewClient(server.URL, "user", "password", 2*time.Second)
	require.NoError(t, err)

	rType, err := client.RepositoryType("445566")
	require.NoError(t, err)
	assert.Equal(t, "githube", rType)
}

func TestBranches(t *testing.T) {
	resBody := []byte(`{"errors":[],"detail":[{"branches":[{"name":"some-branch-name","url":"https://github.com/your-org/your-repo/tree/some-branch-name","createPullRequestUrl":"https://github.com/your-org/your-repo/compare","repository":{"name":"your-repo","url":"https://github.com/your-org/your-repo"}}]}]}`)
	server := httptest.NewServer(testHandler(resBody))
	defer server.Close()
	client, err := gojira.NewClient(server.URL, "user", "password", 2*time.Second)
	require.NoError(t, err)

	branches, err := client.Branches("112233", "github")
	require.NoError(t, err)
	require.Len(t, branches, 1)
	assert.Equal(t, "some-branch-name", branches[0].Name)
	assert.Equal(t, "https://github.com/your-org/your-repo/tree/some-branch-name", branches[0].URL)
	assert.Equal(t, "your-repo", branches[0].Repository.Name)
	assert.Equal(t, "https://github.com/your-org/your-repo", branches[0].Repository.URL)
}

func testHandler(resBody []byte) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write(resBody)
	})
}
