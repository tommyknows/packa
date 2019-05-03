package command

import (
	"io/ioutil"
	"math/rand"
	"os"
	"os/exec"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInstall(t *testing.T) {
	cmdH, err := NewHandler()
	assert.Nil(t, err)
	cmd := "-h"
	expectedOutput := `usage: go get [-d] [-m] [-u] [-v] [-insecure] [build flags] [packages]
Run 'go help get' for details.
`
	output, err := cmdH.Install(cmd)
	assert.NotNil(t, err)
	assert.Equal(t, expectedOutput, output)
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
	_, err = cmdH.Install("github.com/thockin/test")
	assert.Nil(t, err)

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
