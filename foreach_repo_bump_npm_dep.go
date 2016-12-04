package main

import (
	"github.com/transcovo/foreachrepo/repos_explorator"
	"net/http"
	"github.com/transcovo/foreachrepo/git"
	"log"
	"os"
	"flag"
	"github.com/transcovo/foreachrepo/npm"
)

const EXAMPLE = `. Example:

foreach_repo_bump_npm_dep -org transcovo -npm-dep chpr-metric -npm-dep-ver 1.0.0` +
	` -branch fixed-chpr-metrics-version -message "TECH Use fixed version for chpr-metric"
`

func main() {
	if !git.IsGitInstalled() {
		log.Fatal("git command not found")
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

	repos, err := repos_explorator.GetReposList(http.DefaultClient, *organization)
	if err != nil {
		panic(err)
	}
	for _, repo := range repos {
		bump_npm_dependency(repo, *npmDep, *npmDepVersion, *branchName, *commitMessage)
	}
}

func bump_npm_dependency(repo repos_explorator.Repo, npmDep string, npmDepVersion string, branchName string, commitMessage string) error {
	dir, err := git.Clone(repo.GitUrl)
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)

	err = npm.UpdatePackage(dir, npmDep, npmDepVersion)
	if err == nil {
		err = git.CommitAndPushInNewBranch(dir, branchName, commitMessage)
		if err == nil {
			log.Println(repo.Name, " -> done!")
		} else {
			log.Println(repo.Name, " -> failed: ", err.Error())
		}

	} else {
		log.Println(repo.Name, " -> skip: ", err.Error())
	}

	return nil
}
