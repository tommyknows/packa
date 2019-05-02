package packages

import (
	"testing"

	"git.ramonruettimann.ml/ramon/packago/app/apis/config"
	"github.com/stretchr/testify/assert"
)

func TestCreatePackage(t *testing.T) {
	tests := []struct {
		url string
		pkg Package
	}{
		{
			url: "test.com/package@latest",
			pkg: Package{
				&config.Package{
					URL:     "test.com/package",
					Version: "latest",
				},
				nil,
			},
		},
		{
			url: "abc.def/another/subpackage@v0.0.1",
			pkg: Package{
				&config.Package{
					URL:     "abc.def/another/subpackage",
					Version: "v0.0.1",
				},
				nil,
			},
		},
		{
			url: "abc.def/nogiven/version",
			pkg: Package{
				&config.Package{
					URL:     "abc.def/nogiven/version",
					Version: "latest",
				},
				nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.url, func(t *testing.T) {
			p := NewPackage(tt.url, nil)
			assert.Equal(t, tt.pkg.URL, p.URL)
			assert.Equal(t, tt.pkg.Version, p.Version)
		})
	}
}
