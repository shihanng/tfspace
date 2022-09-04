package space_test

import (
	"testing"

	"github.com/shihanng/tfspace/space"
	"gotest.tools/v3/assert"
)

func TestAddBackend(t *testing.T) {
	t.Parallel()

	var testSpaces space.Spaces

	testSpaces.AddBackend("dev", "backend.dev")
	testSpaces.AddBackend("dev", "be.dev")
	testSpaces.AddBackend("dev", "backend.dev")

	expected := space.Spaces{
		{
			Name:    "dev",
			Backend: []string{"backend.dev", "be.dev"},
		},
	}

	assert.DeepEqual(t, testSpaces, expected)
}

func TestRemoveBackend(t *testing.T) {
	t.Parallel()

	testSpaces := space.Spaces{
		{
			Name:    "dev",
			Backend: []string{"backend.dev", "be.dev"},
		},
	}

	testSpaces.RemoveBackend("dev", "backend.dev")
	testSpaces.RemoveBackend("stg", "backend.stg")

	assert.DeepEqual(t, testSpaces, space.Spaces{
		{
			Name:    "dev",
			Backend: []string{"be.dev"},
		},
	})
}
