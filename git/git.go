package git

import (
	"io/ioutil"
	"os/exec"
)

func IsGitInstalled() bool {
	cmd := exec.Command("git", "--version")
	clone_start_err := cmd.Start()
	if clone_start_err != nil {
		return false
	}
	clone_end_err := cmd.Wait()
	if clone_end_err != nil {
		return false
	}
	return true
}

func CloneRepo(url string) (string, error) {
	dir, mkdir_err := ioutil.TempDir("", "")
	if mkdir_err != nil {
		return "", mkdir_err
	}

	cmd := exec.Command("git", "clone", url, dir)
	clone_start_err := cmd.Start()
	if clone_start_err != nil {
		return "", clone_start_err
	}
	clone_end_err := cmd.Wait()
	if clone_end_err != nil {
		return "", clone_end_err
	}
	return dir, nil
}
