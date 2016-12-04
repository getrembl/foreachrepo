package git

import (
	"testing"
	"io/ioutil"
	"os"
	"path/filepath"
	"github.com/stretchr/testify/assert"
	"strings"
	"sort"
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
	executor := &Executor{dir}
	executor.Exec("git", "init")
	executor.Exec("git", "init", "--bare", origin)
	executor.Exec("git", "remote", "add", "origin", origin)

	executor.Exec("touch", "file1.txt")
	executor.Exec("git", "add", ".")
	executor.Exec("git", "commit", "-m", `"Initial commit"`)
	executor.Exec("git", "push", "-u", "origin", "master")

	executor.Exec("git", "checkout", "-b", "a-branch")
	executor.Exec("touch", "file2.txt")
	executor.Exec("git", "add", ".")
	executor.Exec("git", "commit", "-m", `"Add a file"`)
	executor.Exec("git", "push", "-u", "origin", "a-branch")

	executor.Exec("git", "checkout", "master")

	return
}

func TestClone(t *testing.T) {
	dir, origin := makeRepo()
	defer os.RemoveAll(dir)
	defer os.RemoveAll(origin)

	cloneDir, cloneErr := Clone(origin)
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

func TestCommitAndPushInNewBranch(t *testing.T) {
	dir, origin := makeRepo()
	defer os.RemoveAll(dir)
	defer os.RemoveAll(origin)

	cloneDir, cloneErr := Clone(origin)
	if cloneErr != nil {
		panic(cloneErr)
	}
	defer os.RemoveAll(cloneDir)

	ExecCommand(cloneDir, "touch", "file3.txt")
	CommitAndPushInNewBranch(cloneDir, "add-3", "Add an other file")

	cloneDirCheck, cloneErrCheck := Clone(origin)
	if cloneErrCheck != nil {
		panic(cloneErrCheck)
	}
	defer os.RemoveAll(cloneDirCheck)

	ExecCommand(cloneDirCheck, "git", "checkout", "add-3")
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
