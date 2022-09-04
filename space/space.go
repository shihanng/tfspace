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

func spaceHasName(name string) func(*Space) bool {
	return func(space *Space) bool {
		return space.Name == name
	}
}
