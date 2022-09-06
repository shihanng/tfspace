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
