package main

import (
	"github.com/transcovo/foreachrepo/repos_explorator"
	"net/http"
	"github.com/transcovo/foreachrepo/git"
	"log"
	"os"
)

func main() {
	if !git.IsGitInstalled() {
		log.Fatal("git command not found")
	}

	repos, err := repos_explorator.GetReposList(http.DefaultClient, "transcovo")
	if err != nil {
		panic(err)
	}
	for _, repo := range repos {
		bump_npm_dependency(repo.GitUrl)
	}
}

func bump_npm_dependency(url string) error {
	dir, err := git.CloneRepo(url)
	if err != nil {
		return err
	}
	log.Print(dir)
	defer os.RemoveAll(dir)
	return nil
}
