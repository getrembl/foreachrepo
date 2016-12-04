package main

import (
	"github.com/transcovo/foreachrepo/git"
	"log"
	"os"
	"flag"
	"github.com/transcovo/foreachrepo/npm"
	"github.com/transcovo/foreachrepo/github"
	"strings"
)

const EXAMPLE = `. Example:

foreachrepo -org transcovo -npm-dep chpr-metrics -npm-dep-ver 1.0.0` +
	` -branch fixed-chpr-metrics-version -message "TECH Use fixed version for chpr-metrics"
`

func main() {
	if !git.IsGitInstalled() {
		log.Fatal("git command not found")
	}

	githubUsername := os.Getenv("GITHUB_USERNAME")
	if githubUsername == "" {
		log.Fatalln("Missing environement variable GITHUB_USERNAME")
	}

	githubPassword := os.Getenv("GITHUB_PASSWORD")
	if githubPassword == "" {
		log.Fatalln("Missing environement variable GITHUB_PASSWORD")
	}

	organization := flag.String("org", "REQUIRED", "The organization to scan")
	npmDep := flag.String("npm-dep", "REQUIRED", "The npm dependency to update")
	npmDepVersion := flag.String("npm-dep-ver", "REQUIRED", "The new version to apply everywhere")
	branchName := flag.String("branch", "REQUIRED", "The branch name to use")
	commitMessage := flag.String("message", "REQUIRED", "The commit message to use")

	flag.Parse()

	if *organization == "REQUIRED" {
		log.Fatalln("organization flag required", EXAMPLE)
	}
	if *npmDep == "REQUIRED" {
		log.Fatalln("npmDep flag required", EXAMPLE)
	}
	if *npmDepVersion == "REQUIRED" {
		log.Fatalln("npmDepVersion flag required", EXAMPLE)
	}
	if *branchName == "REQUIRED" {
		log.Fatalln("branchName flag required", EXAMPLE)
	}
	if *commitMessage == "REQUIRED" {
		log.Fatalln("commitMessage flag required", EXAMPLE)
	}

	httpInterface := &github.AuthHttpInterface{githubUsername, githubPassword}

	repos, err := github.GetReposList(httpInterface, *organization)
	if err != nil {
		panic(err)
	}
	urls := []string{}
	for _, repo := range repos {
		url := bump_npm_dependency(httpInterface, repo, *npmDep, *npmDepVersion, *branchName, *commitMessage)
		if url != "" {
			urls = append(urls, url)
		}
	}
	println("===== Done =====")
	println(strings.Join(urls, "\n"))
}

func bump_npm_dependency(httpInterface *github.AuthHttpInterface, repo github.Repo, npmDep string, npmDepVersion string,
branchName string, commitMessage string) string {
	dir, err := git.Clone(repo.GitUrl)
	if err != nil {
		log.Println(repo.Name, " -> failed: ", err.Error())
		return ""
	}
	defer os.RemoveAll(dir)

	err = npm.UpdatePackage(dir, npmDep, npmDepVersion)
	if err == nil {
		err = git.CommitAndPushInNewBranch(dir, branchName, commitMessage)
		if err == nil {
			url := github.CreatePullRequest(httpInterface, repo, branchName, commitMessage)

			log.Println(repo.Name, " -> done! (", url, ")")
			return url
		}

		log.Println(repo.Name, " -> failed: ", err.Error())
		return ""
	}

	log.Println(repo.Name, " -> skip: ", err.Error())
	return ""
}
