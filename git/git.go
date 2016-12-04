package git

import (
	"io/ioutil"
	"os/exec"
	"log"
)

type Executor struct {
	dir string
}

func (E *Executor) Exec(name string, elements ...string) error {
	cmd := exec.Command(name, elements...)
	if E.dir != "" {
		cmd.Dir = E.dir
	}
	start_err := cmd.Start()
	if start_err != nil {
		return start_err
	}
	wait_err := cmd.Wait()
	if wait_err != nil {
		return wait_err
	}
	return nil
}

func ExecCommand(dir string, name string, elements ...string) error {
	executor := &Executor{dir}
	return executor.Exec(name, elements...)
}

func IsGitInstalled() bool {
	return ExecCommand("", "git", "--version") == nil
}

func Clone(url string) (string, error) {
	log.Println("Cloning repo ", url)
	dir, mkdirErr := ioutil.TempDir("", "")
	if mkdirErr != nil {
		return "", mkdirErr
	}

	gitErr := ExecCommand(dir, "git", "clone", url, dir)
	if gitErr != nil {
		return "", gitErr
	}
	return dir, nil
}

func CommitAndPushInNewBranch(dir string, branch string, message string) error {
	executor := &Executor{dir}
	err := executor.Exec("git", "checkout", "-b", branch)
	if err == nil {
		err = executor.Exec("git", "add", ".")
	}
	if err == nil {
		err = executor.Exec("git", "commit", "-m", message)
	}
	if err == nil {
		err = executor.Exec("git", "push", "-u", "origin", branch)
	}
	return err
}
