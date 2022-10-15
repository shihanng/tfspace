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

func TestEnv(t *testing.T) {
	t.Parallel()

	testSpaces := space.Spaces{
		{
			Name:      "dev",
			Backend:   []string{"backend.dev", "be.dev"},
			Varfile:   []string{"abc.tfvars"},
			Workspace: "dev_ws",
		},
		{
			Name:      "stg",
			Backend:   []string{"backend.stg", "b.stg"},
			Varfile:   []string{"stg.tfvars", "secret.tfvars"},
			Workspace: "stg",
		},
	}

	tests := []struct {
		name        string
		expected    []string
		expectedErr interface{}
	}{
		{
			name: "dev",
			expected: []string{
				`TFSPACE=dev`,
				`TF_CLI_ARGS_init='-backend-config="backend.dev" -backend-config="be.dev"'`,
				`TF_CLI_ARGS_plan='-var-file="abc.tfvars"'`,
				`TF_CLI_ARGS_apply='-var-file="abc.tfvars"'`,
				`TF_WORKSPACE=dev_ws`,
			},
			expectedErr: noError,
		},
		{
			name: "stg",
			expected: []string{
				`TFSPACE=stg`,
				`TF_CLI_ARGS_init='-backend-config="backend.stg" -backend-config="b.stg"'`,
				`TF_CLI_ARGS_plan='-var-file="stg.tfvars" -var-file="secret.tfvars"'`,
				`TF_CLI_ARGS_apply='-var-file="stg.tfvars" -var-file="secret.tfvars"'`,
				"TF_WORKSPACE=stg",
			},
			expectedErr: noError,
		},
		{
			name:        "prod",
			expected:    nil,
			expectedErr: space.ErrNotFound,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			actual, err := testSpaces.Env(tt.name)
			assert.DeepEqual(t, actual, tt.expected)
			assert.ErrorType(t, err, tt.expectedErr)
		})
	}
}

func noError(err error) bool {
	return err == nil
}
