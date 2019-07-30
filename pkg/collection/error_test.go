package collection

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAll(t *testing.T) {
	var e1, e2 Error
	e2.Add("test", fmt.Errorf("test"))
	e2.Add("test2", fmt.Errorf("test2"))
	assert.Nil(t, e1.IfNotEmpty())
	assert.NotNil(t, e2.IfNotEmpty())
	assert.Equal(t, "test2", e2.errors["test2"].Error())

	e1.Merge(e2)
	assert.NotNil(t, e1.IfNotEmpty())
}
