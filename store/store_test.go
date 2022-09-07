package store_test

import (
	"testing"

	"github.com/shihanng/tfspace/space"
	"github.com/shihanng/tfspace/store"
	"gotest.tools/v3/assert"
)

func TestLoad(t *testing.T) {
	actual, err := store.Load("./testdata/tfspace.yml")

	expected := space.Spaces{
		{
			Name:      "dev",
			Backend:   []string{"dev.backend"},
			Varfile:   []string{"dev.tfvars"},
			Workspace: "dev",
		},
		{
			Name:    "stg",
			Backend: []string{"stg.backend", "stg.be"},
			Varfile: []string{"stg.tfvars", "stg-secret.tfvars"},
		},
	}

	assert.NilError(t, err)
	assert.DeepEqual(t, actual, expected)
}

func TestLoad_Errors(t *testing.T) {
	tests := []struct {
		path   string
		errMsg string
	}{
		{"./testdata/load_bad_yml.yml", "decode mapstructure"},
		{"./testdata/load_int_key.yml", "name is not string"},
		{"./testdata/load_not_hash.yml", "value is not hash map"},
		{"./testdata/load_text.yml", "decode yaml"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.path, func(t *testing.T) {
			t.Parallel()

			actual, err := store.Load(tt.path)
			assert.ErrorContains(t, err, tt.errMsg)
			assert.Assert(t, actual == nil)
		})
	}
}
