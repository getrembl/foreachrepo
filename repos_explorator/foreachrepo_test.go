package repos_explorator

import (
	"testing"
	"net/http"
	"errors"
	"bytes"
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
		"ssh_url": "git@github.com:org/repo1.git"
	}, {
		"name": "repo2",
		"ssh_url": "git@github.com:org/repo2.git"
	}]`

	page2 := `[{
		"name": "repo3",
		"ssh_url": "git@github.com:org/repo3.git"
	}, {
		"name": "repo4",
		"ssh_url": "git@github.com:org/repo4.git"
	}]`

	page3 := "[]"

	responses := map[string]string{
		"https://api.github.com/orgs/org/repos?page=1&per_page=50": page1,
		"https://api.github.com/orgs/org/repos?page=2&per_page=50": page2,
		"https://api.github.com/orgs/org/repos?page=3&per_page=50": page3,
	}

	getter := &TestHttpGetterSuccess{responses: responses}
	repos, err := GetReposList(getter, "org")
	if err != nil {
		t.Error("Expected no err, got err.Error()=", err.Error())
	}
	if len(repos) != 4 {
		t.Error("Expected len(repos) to be 4, got ", len(repos))
	}
	if repos[0].GitUrl != "git@github.com:org/repo1.git" {
		t.Error("Expected repos[0].GitUrl to be git@github.com:org/repo1.git, got ", repos[0].GitUrl)
	}
	if repos[1].GitUrl != "git@github.com:org/repo2.git" {
		t.Error("Expected repos[1].GitUrl to be git@github.com:org/repo2.git, got ", repos[1].GitUrl)
	}
	if repos[2].GitUrl != "git@github.com:org/repo3.git" {
		t.Error("Expected repos[2].GitUrl to be git@github.com:org/repo3.git, got ", repos[2].GitUrl)
	}
	if repos[3].GitUrl != "git@github.com:org/repo4.git" {
		t.Error("Expected repos[3].GitUrl to be git@github.com:org/repo4.git, got ", repos[3].GitUrl)
	}
}
