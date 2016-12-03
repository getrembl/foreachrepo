package foreachrepo

import (
	"net/http"
	"net/url"
	"fmt"
	"encoding/json"
	"strconv"
	"log"
)

type GitUserConfig struct {
	Username string
	Password string
}

type Repo struct {
	Name   string
	GitUrl string
}

type HttpGetter interface {
	Get(url string) (resp *http.Response, err error)
}

type githubApiRepoDescription struct {
	Name    string
	Ssh_url string
}

func getJson(httpGetter HttpGetter, url string, target interface{}) error {
	r, err := httpGetter.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(target)
}

func appendPageRepos(getter HttpGetter, repos *[]Repo, pageUrl *url.URL, i int) (int, error) {
	query := pageUrl.Query()
	query.Set("page", strconv.Itoa(i))
	query.Set("per_page", "50")
	pageUrl.RawQuery = query.Encode()


	log.Print("Loading page ", i)

	page := make([]githubApiRepoDescription,0)
	err := getJson(getter, pageUrl.String(), &page)
	if err != nil {
		return 0, err
	}

	log.Print("Page len = ", len(page))

	for _, repoDescription := range page {
		log.Print("processing repo ", repoDescription.Name)

		repo := Repo{
			Name:repoDescription.Name,
			GitUrl:repoDescription.Ssh_url,
		}

		*repos = append(*repos, repo)
	}
	return len(page), nil
}

func GetReposList(getter HttpGetter, organization string) ([]Repo, error) {
	pageUrl := &url.URL{
		Scheme:"https",
		Host:"api.github.com",
		Path:fmt.Sprintf("orgs/%v/repos", organization),
	}
	repos := []Repo{}
	for i := 1; ; i++ {
		count, err := appendPageRepos(getter, &repos, pageUrl, i)
		if err != nil {
			return nil, err
		}
		if count == 0 {
			return repos, nil
		}
	}
}
