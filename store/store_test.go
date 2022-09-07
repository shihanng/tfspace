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

func TestLoad_BadYML(t *testing.T) {
	actual, err := store.Load("./testdata/load_bad_yml.yml")

	assert.ErrorContains(t, err, "source data must be an array or slice")
	assert.Assert(t, actual == nil)
}

func TestLoad_IntKey(t *testing.T) {
	actual, err := store.Load("./testdata/load_int_key.yml")

	assert.ErrorContains(t, err, "name is not string")
	assert.Assert(t, actual == nil)
}

func TestLoad_NotHash(t *testing.T) {
	actual, err := store.Load("./testdata/load_not_hash.yml")

	assert.ErrorContains(t, err, "value is not hash map")
	assert.Assert(t, actual == nil)
}

func TestLoad_Text(t *testing.T) {
	actual, err := store.Load("./testdata/load_text.yml")

	assert.ErrorContains(t, err, "string was used where mapping")
	assert.Assert(t, actual == nil)
}
