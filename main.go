package main

import (
	"github.com/transcovo/foreachrepo/git"
	"log"
	"os"
	"flag"
	"github.com/transcovo/foreachrepo/npm"
	"github.com/transcovo/foreachrepo/github"
	"strings"
	"github.com/transcovo/foreachrepo/tasks"
)

const EXAMPLES = `. Examples:

Bump a single npm dependency to a specific version in all repos:

$> foreachrepo -task BUMP -org transcovo -npm-dep chpr-metrics -npm-dep-ver 1.0.0` +
	` -branch fixed-chpr-metrics-version -message "TECH Use fixed version for chpr-metrics"

Freeze all package.json dependencies to the current exact version of the current result of npm intall

$> foreachrepo -task FREEZE -org transcovo` +
	` -branch freeze-all-deps -message "TECH Freeze all dependencies to the current result of npm i"
`

func main() {
	g := git.Git("")
	if !g.IsInstalled() {
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

	// generic, mandatory
	taskName := flag.String("task", "DEFAULT", "The task to execute")
	organization := flag.String("org", "DEFAULT", "The organization to scan")
	branchName := flag.String("branch", "DEFAULT", "The branch name to use")
	commitMessage := flag.String("message", "DEFAULT", "The commit message to use")

	// for bumping single dependency parameter
	npmDep := flag.String("npm-dep", "DEFAULT", "The npm dependency to update")
	npmDepVersion := flag.String("npm-dep-ver", "DEFAULT", "The new version to apply everywhere")

	flag.Parse()

	if *organization == "DEFAULT" {
		log.Fatalln("organization flag required", EXAMPLES)
	}
	if *taskName == "DEFAULT" {
		log.Fatalln("task flag required", EXAMPLES)
	}
	if *branchName == "DEFAULT" {
		log.Fatalln("branch-name flag required", EXAMPLES)
	}
	if *commitMessage == "DEFAULT" {
		log.Fatalln("commit-message flag required", EXAMPLES)
	}

	var task tasks.Task

	if *taskName == "BUMP" {
		if *npmDep == "DEFAULT" {
			log.Fatalln("npm-dep flag required when task is BUMP", EXAMPLES)
		}
		if *npmDepVersion == "DEFAULT" {
			log.Fatalln("npm-dep-version flag required when task is BUMP", EXAMPLES)
		}
		task = BumpNpmDependencyTask{
			npmDep:*npmDep,
			npmDepVersion:*npmDepVersion,
		}
	} else if *taskName == "FREEZE" {
		task = FreezeTask{}
	} else {
		log.Fatalln("Unknown task type ", *taskName, EXAMPLES)
	}

	httpInterface := &github.AuthHttpInterface{githubUsername, githubPassword}

	repos, err := github.GetReposList(httpInterface, *organization)
	if err != nil {
		panic(err)
	}
	urls := []string{}
	for _, repo := range repos {
		url := tasks.ExecuteTask(httpInterface, repo, task, *branchName, *commitMessage)
		if url != "" {
			urls = append(urls, url)
		}
	}
	println("===== Done =====")
	println(strings.Join(urls, "\n"))
}

type BumpNpmDependencyTask struct {
	npmDep        string
	npmDepVersion string
}

func (t BumpNpmDependencyTask) Execute(dir string) error {
	return npm.UpdatePackage(dir, t.npmDep, t.npmDepVersion)
}

type FreezeTask struct{}

func (t FreezeTask) Execute(dir string) error {
	return npm.FreezePackage(dir)
}
