package command

import (
	"io/ioutil"
	"math/rand"
	"os"
	"os/exec"
	"path"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestInstall(t *testing.T) {
	expectedOutput := "go get test@: invalid module version syntax\n"
	expectedErrMessage := expectedOutput + "exit status 1"
	// create temporary directory
	tmpWorkDir, err := ioutil.TempDir("", "packa-test")
	assert.Nil(t, err)
	defer os.RemoveAll(tmpWorkDir)
	// create the handler
	cmdH, err := NewHandler(WorkingDir(tmpWorkDir))
	assert.Nil(t, err)

	cmd := "test"
	output, err := cmdH.Install(cmd, "")
	assert.Equal(t, "", output)
	_, ok := err.(InstallError)
	assert.True(t, ok)

	_, ok = errors.Cause(err).(*exec.ExitError)
	assert.True(t, ok)

	assert.Equal(t, expectedErrMessage, err.Error())
}

func TestWorkingDir(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping working dir as network operation in short mode")
	}

	// setup
	dir, err := ioutil.TempDir("", "test")
	defer os.RemoveAll(dir)
	assert.Nil(t, err)
	// by having a go mod file in this directory,
	// we see if this is actually working by checking
	// the modules file afterwards
	initMod := exec.Command("go", "mod", "init", "test")
	initMod.Dir = dir
	_, err = initMod.CombinedOutput()
	assert.Nil(t, err)

	cmdH, err := NewHandler(WorkingDir(dir))
	assert.Nil(t, err)
	_, err = cmdH.Install("github.com/thockin/test", "latest")
	assert.Nil(t, err)

	// check that the go mod file contains the git repo
	mod, err := ioutil.ReadFile(path.Join(dir, "go.mod"))
	assert.Nil(t, err)
	assert.Contains(t, string(mod), "github.com/thockin/test")
}

func TestRemove(t *testing.T) {
	cmdH, err := NewHandler()
	assert.Nil(t, err)
	// generate a random string of length 20
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, 20)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}

	err = cmdH.Remove(string(b))
	// RemoveAll should never return an error
	assert.Nil(t, err)
}

func TestGetVersion(t *testing.T) {
	cmdH, err := NewHandler()
	assert.Nil(t, err)
	tests := []struct {
		url     string
		output  string
		version string
	}{
		{
			url:     "test.com/no/bla",
			output:  "go: extracting test.com/no/bla v0.0.1\n",
			version: "v0.0.1",
		},
		{
			url:     "test.com/no/bla",
			output:  "go: getting someth.in/else v0.0.1\ngo: extracting someth.in/else v0.0.1\n",
			version: "",
		},
		{
			url:     "test.com/no/bla",
			output:  "multiline\ngo: extracting test.com/no/bla v0.0.1\n",
			version: "v0.0.1",
		},
		{
			url:     "test.com/no/bla",
			output:  "invalid output\n",
			version: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.output, func(t *testing.T) {
			v := cmdH.getVersion(tt.url, tt.output)
			assert.Equal(t, tt.version, v)
		})
	}
}
