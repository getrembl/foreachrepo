package git

import (
	"io/ioutil"
	"os/exec"
	"log"
	"errors"
)

type Sys interface {
	Command(name string, arg ...string) *exec.Cmd
	TempDir(dir, prefix string) (name string, err error)
}

type ActualSys struct{}

func (c ActualSys) Command(name string, arg ...string) *exec.Cmd {
	return exec.Command(name, arg...)
}
func (c ActualSys) TempDir(dir, prefix string) (string, error) {
	return ioutil.TempDir(dir, prefix)
}

type git struct {
	Sys Sys
	Dir string
}

func Git(dir string) *git {
	return &git{&ActualSys{}, dir}
}

func (g *git) setSys(sys Sys) {
	g.Sys = sys
}

func (g *git) Clone(url string) (string, error) {
	if g.Dir != "" {
		return "", errors.New("This git repo's dir is already initialized: " + g.Dir)
	}
	log.Println("Cloning repo ", url)
	dir, mkdirErr := g.Sys.TempDir("", "")
	if mkdirErr != nil {
		return "", mkdirErr
	}
	g.Dir = dir

	gitErr := g.Exec("git", "clone", url, dir)
	if gitErr != nil {
		return "", gitErr
	}
	return dir, nil
}

func (g *git) Exec(name string, elements ...string) error {
	cmd := g.Sys.Command(name, elements...)
	if g.Dir != "" {
		cmd.Dir = g.Dir
	}
	startErr := cmd.Start()
	if startErr != nil {
		return startErr
	}
	waitErr := cmd.Wait()
	if waitErr != nil {
		return waitErr
	}
	return nil
}

func (g *git) IsInstalled() bool {
	return g.Exec("git", "--version") == nil
}

func (g *git) CommitAndPushInNewBranch(branch string, message string) error {
	err := g.Exec("git", "checkout", "-b", branch)
	if err == nil {
		err = g.Exec("git", "add", ".")
	}
	if err == nil {
		err = g.Exec("git", "commit", "-m", message)
	}
	if err == nil {
		err = g.Exec("git", "push", "-u", "origin", branch)
	}
	return err
}
