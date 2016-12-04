package github

import (
	"testing"
	"net/http"
	"errors"
	"bytes"
	"github.com/stretchr/testify/assert"
)

type TestHttpGetterError struct{}

func (getter TestHttpGetterError) Get(url string) (resp *http.Response, err error) {
	return nil, errors.New("Mock error")
}

type ClosingBuffer struct {
	*bytes.Buffer
}

func (cb *ClosingBuffer) Close() error {
	//we don't actually have to do anything here, since the buffer is
	// just some data in memory
	//and the error is initialized to no-error
	return nil
}

type TestHttpGetterSuccess struct {
	responses map[string]string
}

func (getter TestHttpGetterSuccess) Get(url string) (*http.Response, error) {
	response := &http.Response{}

	responseString, ok := getter.responses[url]
	if !ok {
		panic(errors.New("Unexpected url: " + url))
	}

	response.Body = &ClosingBuffer{bytes.NewBufferString(responseString)}
	return response, nil
}

func TestGetReposError(t *testing.T) {
	getter := &TestHttpGetterError{}
	repos, err := GetReposList(getter, "org")
	if repos != nil {
		t.Error("When the HTTP call fails, returned repos should be nil, repos=", repos)
	}
	if err.Error() != "Mock error" {
		t.Error("When the HTTP call fails, the original error should be propagated, err.Error()=", err.Error())
	}
}
func TestGetReposSuccess(t *testing.T) {
	page1 := `[{
		"name": "repo1",
		"ssh_url": "git@github.com:org/repo1.git",
		"pulls_url": "http://api.github.com/repos/org/repo1/pulls{/number}"
	}, {
		"name": "repo2",
		"ssh_url": "git@github.com:org/repo2.git",
		"pulls_url": "http://api.github.com/repos/org/repo2/pulls{/number}"
	}]`

	page2 := `[{
		"name": "repo3",
		"ssh_url": "git@github.com:org/repo3.git",
		"pulls_url": "http://api.github.com/repos/org/repo3/pulls{/number}"
	}, {
		"name": "repo4",
		"ssh_url": "git@github.com:org/repo4.git",
		"pulls_url": "http://api.github.com/repos/org/repo4/pulls{/number}"
	}]`

	page3 := "[]"

	responses := map[string]string{
		"https://api.github.com/orgs/org/repos?page=1&per_page=50": page1,
		"https://api.github.com/orgs/org/repos?page=2&per_page=50": page2,
		"https://api.github.com/orgs/org/repos?page=3&per_page=50": page3,
	}

	getter := &TestHttpGetterSuccess{responses: responses}
	repos, err := GetReposList(getter, "org")
	assert.Nil(t, err)
	assert.Len(t, repos, 4)
	assert.Equal(t, "git@github.com:org/repo1.git", repos[0].GitUrl)
	assert.Equal(t, "repo1", repos[0].Name)
	assert.Equal(t, "http://api.github.com/repos/org/repo1/pulls", repos[0].PullsUrl)
	
	assert.Equal(t, "git@github.com:org/repo2.git", repos[1].GitUrl)
	assert.Equal(t, "repo2", repos[1].Name)
	assert.Equal(t, "http://api.github.com/repos/org/repo2/pulls", repos[1].PullsUrl)
	
	assert.Equal(t, "git@github.com:org/repo3.git", repos[2].GitUrl)
	assert.Equal(t, "repo3", repos[2].Name)
	assert.Equal(t, "http://api.github.com/repos/org/repo3/pulls", repos[2].PullsUrl)
	
	assert.Equal(t, "git@github.com:org/repo4.git", repos[3].GitUrl)
	assert.Equal(t, "repo4", repos[3].Name)
	assert.Equal(t, "http://api.github.com/repos/org/repo4/pulls", repos[3].PullsUrl)
}
