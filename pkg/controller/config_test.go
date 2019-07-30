package controller

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestConfigIO creates a source and a destination file,
// writes a sample config into source, tests the parsing
// of the config and writes it to destFile
func TestConfigIO(t *testing.T) {
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
	assert.NoError(t, err)
	destFile, err = ioutil.TempFile("", "")
	assert.NoError(t, err)
	sourceFile.Close()
	destFile.Close()
	defer os.Remove(sourceFile.Name())
	defer os.Remove(destFile.Name())

	// write the config into the source file
	assert.NoError(t, ioutil.WriteFile(sourceFile.Name(), testCfg, 0700))

	opt := ConfigFile(sourceFile.Name())
	ctl := &Controller{}
	assert.NoError(t, opt(ctl))

	// save the file to the new destination
	ctl.configuration.file = destFile.Name()
	err = ctl.configuration.save()
	assert.NoError(t, err)
	// ensure that the contents are equal
	d, err := ioutil.ReadFile(destFile.Name())
	assert.NoError(t, err)
	assert.Equal(t, testCfg, d)
}
