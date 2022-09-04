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

func TestAddVarfile(t *testing.T) {
	t.Parallel()

	var testSpaces space.Spaces

	testSpaces.AddVarfile("dev", "dev.tfvars")
	testSpaces.AddVarfile("dev", "abc.tfvars")
	testSpaces.AddVarfile("dev", "dev.tfvars")

	expected := space.Spaces{
		{
			Name:    "dev",
			Varfile: []string{"dev.tfvars", "abc.tfvars"},
		},
	}

	assert.DeepEqual(t, testSpaces, expected)
}

func TestRemoveVarfile(t *testing.T) {
	t.Parallel()

	testSpaces := space.Spaces{
		{
			Name:    "dev",
			Varfile: []string{"dev.tfvars", "abc.tfvars"},
		},
	}

	testSpaces.RemoveVarfile("dev", "dev.tfvars")
	testSpaces.RemoveVarfile("stg", "stg.tfvars")

	assert.DeepEqual(t, testSpaces, space.Spaces{
		{
			Name:    "dev",
			Varfile: []string{"abc.tfvars"},
		},
	})
}

func TestSetWorkspace(t *testing.T) {
	t.Parallel()

	var testSpaces space.Spaces

	testSpaces.SetWorkspace("dev", "dev_ws")
	testSpaces.SetWorkspace("dev", "dev_ws_new")

	assert.DeepEqual(t, testSpaces, space.Spaces{
		{
			Name:      "dev",
			Workspace: "dev_ws_new",
		},
	})
}

func TestUnsetWorkspace(t *testing.T) {
	t.Parallel()

	testSpaces := space.Spaces{
		{
			Name:      "dev",
			Workspace: "dev_ws",
		},
	}

	testSpaces.UnsetWorkspace("dev")
	testSpaces.UnsetWorkspace("stg")

	assert.DeepEqual(t, testSpaces, space.Spaces{
		{
			Name: "dev",
		},
	})
}
