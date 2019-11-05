package controller

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"

	"github.com/matryer/is"
)

// TestConfigIO creates a source and a destination file,
// writes a sample config into source, tests the parsing
// of the config and writes it to destFile
func TestConfigIO(t *testing.T) {
	is := is.New(t)

	testCfg := []byte(`packages:
  fakeHandler:
  - url: github.com/test/bla
  thirdHandler:
  - totest: something
settings:
  handler:
    fakeHandler:
      workingDir: /tmp
`)

	// create the source files
	var sourceFile, destFile *os.File
	sourceFile, err := ioutil.TempFile("", "")
	is.NoErr(err)
	destFile, err = ioutil.TempFile("", "")
	is.NoErr(err)
	sourceFile.Close()
	destFile.Close()
	defer os.Remove(sourceFile.Name())
	defer os.Remove(destFile.Name())

	// write the config into the source file
	err = ioutil.WriteFile(sourceFile.Name(), testCfg, 0700)
	is.NoErr(err)

	option := ConfigFile(sourceFile.Name())
	ctl := &Controller{}
	is.NoErr(option(ctl))

	// save the file to the new destination
	ctl.configuration.file = destFile.Name()
	err = ctl.configuration.save()
	is.NoErr(err)
	// ensure that the contents are equal
	d, err := ioutil.ReadFile(destFile.Name())
	is.NoErr(err)
	is.Equal(string(testCfg), string(d))

	// simulate removal of a package.
	// tests for bug in #6
	noPackages := json.RawMessage("[]")
	ctl.configuration.Packages["fakeHandler"] = &noPackages

	err = ctl.configuration.save()
	is.NoErr(err)

	smallerCfg := []byte(`packages:
  fakeHandler: []
  thirdHandler:
  - totest: something
settings:
  handler:
    fakeHandler:
      workingDir: /tmp
`)
	d, err = ioutil.ReadFile(destFile.Name())
	is.NoErr(err)
	is.Equal(string(smallerCfg), string(d))
}
