// Package space manages the settings of space.
// A space is an environment of which Terraform operates on.
package space

import (
	"fmt"
	"strings"

	"github.com/cockroachdb/errors"
	"github.com/samber/lo"
)

var ErrNotFound = errors.New("space not found")

// Space contains the configuration of each space.
type Space struct {
	Name string

	Backend   []string
	Varfile   []string
	Workspace string
}

func (s *Space) addBackend(backend string) {
	index := lo.IndexOf(s.Backend, backend)

	if index >= 0 {
		return
	}

	s.Backend = append(s.Backend, backend)
}

func (s *Space) removeBackend(backend string) {
	rejected := lo.Reject(s.Backend, func(val string, _ int) bool {
		return val == backend
	})
	s.Backend = rejected
}

func (s *Space) addVarfile(varfile string) {
	index := lo.IndexOf(s.Varfile, varfile)

	if index >= 0 {
		return
	}

	s.Varfile = append(s.Varfile, varfile)
}

func (s *Space) removeVarfile(varfile string) {
	rejected := lo.Reject(s.Varfile, func(val string, _ int) bool {
		return val == varfile
	})
	s.Varfile = rejected
}

// Spaces is a list of Space.
type Spaces []*Space

// AddBackend adds backend into the space of name.
// If the space does not exist, a new one will be created.
// If the backend already exists in the space, it will not do anything.
func (s *Spaces) AddBackend(name, backend string) {
	space, found := findSpace(*s, name)

	space.addBackend(backend)

	if !found {
		*s = append(*s, space)
	}
}

// AddVarfile adds var-file into the space of name.
// If the space does not exist, a new one will be created.
// If the var-file already exists in the space, it will not do anything.
func (s *Spaces) AddVarfile(name, varfile string) {
	space, found := findSpace(*s, name)

	space.addVarfile(varfile)

	if !found {
		*s = append(*s, space)
	}
}

// RemoveVarfile removes var-file from the space.
// If the space does not exist, it does not do anything.
// If the var-file does not exist in the space, it will not do anything.
func (s *Spaces) RemoveVarfile(name, varfile string) {
	space, found := findSpace(*s, name)

	if !found {
		return
	}

	space.removeVarfile(varfile)
}

// RemoveBackend removes backend from the space.
// If the space does not exist, it does not do anything.
// If the backend does not exist in the space, it will not do anything.
func (s *Spaces) RemoveBackend(name, backend string) {
	space, found := findSpace(*s, name)

	if !found {
		return
	}

	space.removeBackend(backend)
}

// SetWorkspace set the value of workspace to the input value.
// If space does not exist, if does not do anything.
func (s *Spaces) SetWorkspace(name, workspace string) {
	space, found := findSpace(*s, name)

	space.Workspace = workspace

	if !found {
		*s = append(*s, space)
	}
}

// UnsetWorkspace set the value of workspace to empty string.
// If space does not exist, if does not do anything.
func (s *Spaces) UnsetWorkspace(name string) {
	space, found := findSpace(*s, name)

	if !found {
		return
	}

	space.Workspace = ""
}

// Env return list of environment variables in the form of
// key=value that can be passed to exec.Command.Env.
func (s *Spaces) Env(name string) ([]string, error) {
	space, found := findSpace(*s, name)

	if !found {
		return nil, ErrNotFound
	}

	var envs []string

	envs = append(envs, fmt.Sprintf("TFSPACE=%s", name))

	if len(space.Backend) > 0 {
		backends := make([]string, 0, len(space.Backend))
		for _, b := range space.Backend {
			backends = append(backends, fmt.Sprintf("-backend-config=\"%s\"", b))
		}

		envs = append(envs, fmt.Sprintf("TF_CLI_ARGS_init='%s'", strings.Join(backends, " ")))
	}

	if len(space.Varfile) > 0 {
		varfiles := make([]string, 0, len(space.Varfile))
		for _, v := range space.Varfile {
			varfiles = append(varfiles, fmt.Sprintf("-var-file=\"%s\"", v))
		}

		envs = append(envs, fmt.Sprintf("TF_CLI_ARGS_plan='%s'", strings.Join(varfiles, " ")))
		envs = append(envs, fmt.Sprintf("TF_CLI_ARGS_apply='%s'", strings.Join(varfiles, " ")))
	}

	if space.Workspace != "" {
		envs = append(envs, fmt.Sprintf("TF_WORKSPACE=%s", space.Workspace))
	}

	return envs, nil
}

func findSpace(spaces Spaces, name string) (*Space, bool) {
	space, found := lo.Find(spaces, spaceHasName(name))

	if !found {
		space = &Space{Name: name} //nolint:exhaustruct
	}

	return space, found
}

func spaceHasName(name string) func(*Space) bool {
	return func(space *Space) bool {
		return space.Name == name
	}
}
