// Package space manages the settings of space.
// A space is an environment of which Terraform operates on.
package space

import (
	"github.com/samber/lo"
)

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
	space, found := lo.Find(*s, spaceHasName(name))

	if !found {
		space = &Space{Name: name} //nolint:exhaustruct
	}

	space.addBackend(backend)

	if !found {
		*s = append(*s, space)
	}
}

// AddVarfile adds var-file into the space of name.
// If the space does not exist, a new one will be created.
// If the var-file already exists in the space, it will not do anything.
func (s *Spaces) AddVarfile(name, varfile string) {
	space, found := lo.Find(*s, spaceHasName(name))

	if !found {
		space = &Space{Name: name} //nolint:exhaustruct
	}

	space.addVarfile(varfile)

	if !found {
		*s = append(*s, space)
	}
}

// RemoveVarfile removes var-file from the space.
// If the space does not exist, it does not do anything.
// If the var-file does not exist in the space, it will not do anything.
func (s *Spaces) RemoveVarfile(name, varfile string) {
	space, found := lo.Find(*s, spaceHasName(name))

	if !found {
		return
	}

	space.removeVarfile(varfile)
}

// RemoveBackend removes backend from the space.
// If the space does not exist, it does not do anything.
// If the backend does not exist in the space, it will not do anything.
func (s *Spaces) RemoveBackend(name, backend string) {
	space, found := lo.Find(*s, spaceHasName(name))

	if !found {
		return
	}

	space.removeBackend(backend)
}

// SetWorkspace set the value of workspace to the input value.
// If space does not exist, if does not do anything.
func (s *Spaces) SetWorkspace(name, workspace string) {
	space, found := lo.Find(*s, spaceHasName(name))

	if !found {
		space = &Space{Name: name} //nolint:exhaustruct
	}

	space.Workspace = workspace

	if !found {
		*s = append(*s, space)
	}
}

// UnsetWorkspace set the value of workspace to empty string.
// If space does not exist, if does not do anything.
func (s *Spaces) UnsetWorkspace(name string) {
	space, found := lo.Find(*s, spaceHasName(name))

	if !found {
		return
	}

	space.Workspace = ""
}

func spaceHasName(name string) func(*Space) bool {
	return func(space *Space) bool {
		return space.Name == name
	}
}
