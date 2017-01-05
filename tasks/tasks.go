package tasks

import (
	"github.com/transcovo/foreachrepo/github"
	"github.com/transcovo/foreachrepo/git"
	"log"
	"os"
)

type Task interface {
	Execute(dir string) error
}

func ExecuteTask(httpInterface *github.AuthHttpInterface, repo github.Repo, task Task, branchName string, commitMessage string) string {
	defer func() {
		if r := recover(); r != nil {
			err, ok := r.(error)
			if ok {
				log.Println(repo.Name, " -> failed: ", err.Error())
			} else {
				log.Println(repo.Name, " -> failed: unknown error")
			}
		}
	}()

	g := git.Git("")
	dir, err := g.Clone(repo.GitUrl)
	if err != nil {
		log.Println(repo.Name, " -> failed: ", err.Error())
		return ""
	}
	defer os.RemoveAll(dir)

	err = task.Execute(dir)

	if err == nil {
		err = g.CommitAndPushInNewBranch(branchName, commitMessage)
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
