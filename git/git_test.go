package git

import (
	"testing"
	"io/ioutil"
	"os"
	"path/filepath"
	"github.com/stretchr/testify/assert"
	"strings"
	"sort"
	"os/exec"
	"errors"
)

func tempDir() string {
	dir, mkdir_err := ioutil.TempDir("", "")
	if mkdir_err != nil {
		panic(mkdir_err)
	}
	return dir
}

func makeRepo() (dir string, origin string) {
	dir = tempDir()
	origin = tempDir()
	g := Git(dir)
	g.Exec("git", "init")
	g.Exec("git", "init", "--bare", origin)
	g.Exec("git", "remote", "add", "origin", origin)

	g.Exec("touch", "file1.txt")
	g.Exec("git", "add", ".")
	g.Exec("git", "commit", "-m", `"Initial commit"`)
	g.Exec("git", "push", "-u", "origin", "master")

	g.Exec("git", "checkout", "-b", "a-branch")
	g.Exec("touch", "file2.txt")
	g.Exec("git", "add", ".")
	g.Exec("git", "commit", "-m", `"Add a file"`)
	g.Exec("git", "push", "-u", "origin", "a-branch")

	g.Exec("git", "checkout", "master")

	return
}

func TestClone(t *testing.T) {
	dir, origin := makeRepo()
	defer os.RemoveAll(dir)
	defer os.RemoveAll(origin)

	g := Git("")
	cloneDir, cloneErr := g.Clone(origin)
	defer os.RemoveAll(cloneDir)

	assert.Nil(t, cloneErr)
	globPattern := filepath.Join(cloneDir, "*.txt")
	files, err := filepath.Glob(globPattern)
	if err != nil {
		panic(err)
	}
	assert.Len(t, files, 1)
	assert.True(t, strings.HasSuffix(files[0], "file1.txt"), "file1.txt mush be the only file present")
}

func TestCloneAlreadyInit(t *testing.T) {
	dir := tempDir()
	defer os.RemoveAll(dir)

	g := Git("/tmp/")
	dir2, err := g.Clone("whatever")
	assert.Empty(t, dir2)
	assert.NotNil(t, err)
}

type TempDirFailSys struct{}

func (c TempDirFailSys) Command(name string, arg ...string) *exec.Cmd {
	return exec.Command(name, arg...)
}
func (c TempDirFailSys) TempDir(dir, prefix string) (string, error) {
	return "", errors.New("Mock error")
}

func TestCloneMkTempDirFail(t *testing.T) {
	dir := tempDir()
	defer os.RemoveAll(dir)

	g := Git("")
	g.setSys(&TempDirFailSys{})
	dir2, err := g.Clone("whatever")
	assert.Empty(t, dir2)
	assert.Equal(t, "Mock error", err.Error())
}

type CommandFailSys struct{}

func (c CommandFailSys) Command(name string, arg ...string) *exec.Cmd {
	return exec.Command("sdjadzejeuqsoieqomequfqomeuzj")
}
func (c CommandFailSys) TempDir(dir, prefix string) (string, error) {
	return ioutil.TempDir(dir, prefix)
}

func TestCloneCommandFail(t *testing.T) {
	dir := tempDir()
	defer os.RemoveAll(dir)

	g := Git("")
	g.setSys(&CommandFailSys{})
	dir2, err := g.Clone("whatever")
	assert.Empty(t, dir2)
	assert.NotNil(t, err)
}

type CommandWaitFailSys struct{}

func (c CommandWaitFailSys) Command(name string, arg ...string) *exec.Cmd {
	return exec.Command("sh", "-c", "sdjadzejeuqsoieqomequfqomeuzj")
}
func (c CommandWaitFailSys) TempDir(dir, prefix string) (string, error) {
	return ioutil.TempDir(dir, prefix)
}
func TestCloneCommandWaitFail(t *testing.T) {
	dir := tempDir()
	defer os.RemoveAll(dir)

	g := Git("")
	g.setSys(&CommandWaitFailSys{})
	dir2, err := g.Clone("whatever")
	assert.Empty(t, dir2)
	assert.NotNil(t, err)
}


func TestIsInstalled(t *testing.T) {
	assert.True(t, Git("").IsInstalled())
}

func TestIsInstalledFail(t *testing.T) {
	g := Git("")
	g.setSys(&CommandFailSys{})
	assert.False(t, g.IsInstalled())
}

func TestCommitAndPushInNewBranch(t *testing.T) {
	dir, origin := makeRepo()
	defer os.RemoveAll(dir)
	defer os.RemoveAll(origin)

	g := Git("")
	cloneDir, cloneErr := g.Clone(origin)
	if cloneErr != nil {
		panic(cloneErr)
	}
	defer os.RemoveAll(cloneDir)

	g.Exec("touch", "file3.txt")
	g.CommitAndPushInNewBranch("add-3", "Add an other file")

	gCheck := Git("")
	cloneDirCheck, cloneErrCheck := gCheck.Clone(origin)
	if cloneErrCheck != nil {
		panic(cloneErrCheck)
	}
	defer os.RemoveAll(cloneDirCheck)

	gCheck.Exec("git", "checkout", "add-3")
	globPattern := filepath.Join(cloneDir, "*.txt")
	files, err := filepath.Glob(globPattern)
	if err != nil {
		panic(err)
	}
	sort.Strings(files)
	assert.Len(t, files, 2)
	assert.True(t, strings.HasSuffix(files[0], "file1.txt"), "file1.txt must be present")
	assert.True(t, strings.HasSuffix(files[1], "file3.txt"), "file3.txt must be present")
}
