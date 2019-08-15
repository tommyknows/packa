package cmd

import (
	"io/ioutil"
	"os"
	"os/user"
	"strings"
	"testing"

	"github.com/matryer/is"
)

func TestGlobalOpts(t *testing.T) {
	is := is.New(t)

	tmpDir, err := ioutil.TempDir("", "")
	is.NoErr(err) // tempdir creation should not fail
	defer os.RemoveAll(tmpDir)
	AddGlobalOptions(WorkingDir(tmpDir))
	defer ResetGlobalOptions()

	out, err := Execute([]string{"sh", "-c", "'pwd'"})
	// cannot directly compare output as on MacOS, TempDir returns /var/...,
	// while the actual directory is reported as /private/var/...
	is.True(strings.HasSuffix(out, tmpDir+"\n")) // command's output should have the dir name and a newline
	is.NoErr(err)                                // executing pwd should not fail
}

func TestExpand(t *testing.T) {
	is := is.New(t)

	usr, err := user.Current()
	is.NoErr(err) // getting current user should not fail

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
			is.Equal(tt.out, out) // actual output and expected output should be the same
			is.NoErr(err)         // we don't expect to receive an error in any test
		})
	}
}

func TestExec(t *testing.T) {
	is := is.New(t)

	out, err := Execute([]string{"echo", "hello world"})
	is.Equal("hello world\n", out) // echo 'hello world' should output hello world
	is.NoErr(err)                  // echo command should not generate an error

	tmpDir, err := ioutil.TempDir("", "")
	is.NoErr(err) // creating tmpdir should not fail
	defer os.RemoveAll(tmpDir)

	out, err = Execute([]string{"sh", "-c", "'pwd'"}, WorkingDir(tmpDir))
	// cannot directly compare output as on MacOS, TempDir returns /var/...,
	// while the actual directory is reported as /private/var/...
	is.True(strings.HasSuffix(out, tmpDir+"\n")) // command's output should have the dir name and a newline
	is.NoErr(err)                                // executing pwd should not fail

	out, err = Execute([]string{"false"})
	is.Equal("", out)   // executing `false` should not print anything
	is.True(err != nil) // error should not be nil

	out, err = Execute([]string{})
	is.Equal("", out)   // not executing any command should output nothing
	is.True(err != nil) // not giving commands should result in an error
}
