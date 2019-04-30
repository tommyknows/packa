package command

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInstall(t *testing.T) {
	cmd := "-h"
	expectedOutput := `usage: go get [-d] [-m] [-u] [-v] [-insecure] [build flags] [packages]
Run 'go help get' for details.
`
	output, err := GoInstall(cmd)
	assert.NotNil(t, err)
	assert.Equal(t, expectedOutput, output)
}
