package cmd

import (
	"io/ioutil"
	"os"
	"os/user"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExpand(t *testing.T) {
	usr, err := user.Current()
	assert.NoError(t, err)
	tests := []struct {
		in  string
		out string
	}{
		{"/tmp", "/tmp"},
		{"~/tests", usr.HomeDir + "/tests"},
		{"~/.packa", usr.HomeDir + "/.packa"},
	}

	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			out, err := expand(tt.in)
			assert.Equal(t, tt.out, out)
			assert.Nil(t, err)
		})
	}
}

func TestExec(t *testing.T) {
	out, err := Execute([]string{"echo", "hello world"})
	assert.Equal(t, "hello world\n", out)
	assert.NoError(t, err)

	tmpDir, err := ioutil.TempDir("", "")
	assert.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	out, err = Execute([]string{"sh", "-c", "'pwd'"}, WorkingDir(tmpDir))
	// cannot directly compare output as on MacOS, TempDir returns /var/...,
	// while the actual directory is reported as /private/var/...
	assert.True(t, strings.HasSuffix(out, tmpDir+"\n"))
	assert.NoError(t, err)

	out, err = Execute([]string{"false"}, DirectPrint(false))
	assert.Equal(t, "", out)
	assert.Error(t, err)

	out, err = Execute([]string{})
	assert.Equal(t, "", out)
	assert.Error(t, err)
}
