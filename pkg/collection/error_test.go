package collection

import (
	"fmt"
	"testing"

	"github.com/matryer/is"
)

func TestAll(t *testing.T) {
	is := is.New(t)
	var e1, e2 Error
	e2.Add("test", fmt.Errorf("test"))
	e2.Add("test2", fmt.Errorf("test2"))
	is.NoErr(e1.IfNotEmpty())
	is.True(e2.IfNotEmpty() != nil)
	is.Equal("test2", e2["test2"].Error())

	e1.Merge(e2)
	is.True(e1.IfNotEmpty() != nil)
}
