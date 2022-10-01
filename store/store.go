// Package store handles load and saving the config file of tfspace.
package store

import (
	"os"
	"path/filepath"

	"github.com/cockroachdb/errors"
	"github.com/goccy/go-yaml"
	"github.com/mitchellh/mapstructure"
	"github.com/shihanng/tfspace/space"
)

// Load config as space.Spaces from path.
func Load(path string) (space.Spaces, error) {
	dat, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		return nil, errors.Wrap(err, "store: read file")
	}

	var v yaml.MapSlice
	if err := yaml.Unmarshal(dat, &v); err != nil {
		return nil, errors.Wrap(err, "store: decode yaml")
	}

	return spacesFromMapSlice(v)
}

// Save config space.Spaces into path.
func Save(path string, spaces space.Spaces) (err error) {
	file, err := os.Create(filepath.Clean(path))
	if err != nil {
		return errors.Wrap(err, "store: create file")
	}

	defer func() {
		err = file.Close()
	}()

	payload := make(yaml.MapSlice, 0, len(spaces))

	for _, space := range spaces {
		if len(space.Backend) == 0 && len(space.Varfile) == 0 && space.Workspace == "" {
			continue
		}

		payload = append(payload, yaml.MapItem{
			Key: space.Name,
			Value: struct {
				Backend   []string `yaml:"backend,omitempty"`
				Varfile   []string `yaml:"varfile,omitempty"`
				Workspace string   `yaml:"workspace,omitempty"`
			}{
				Backend:   space.Backend,
				Varfile:   space.Varfile,
				Workspace: space.Workspace,
			},
		})
	}

	if len(payload) == 0 {
		_, err := file.WriteString("")

		return errors.Wrap(err, "store: write empty yaml to file")
	}

	return errors.Wrap(
		yaml.NewEncoder(file, yaml.IndentSequence(true)).Encode(payload),
		"store: write yaml to file",
	)
}

func spacesFromMapSlice(ms yaml.MapSlice) (space.Spaces, error) {
	spaces := make(space.Spaces, 0, len(ms))

	for _, item := range ms {
		name, ok := item.Key.(string) //nolint:varnamelen
		if !ok {
			return nil, errors.New("store: name is not string")
		}

		values, ok := item.Value.(map[string]interface{})
		if !ok {
			return nil, errors.New("store: value is not hash map")
		}

		space := space.Space{ //nolint:exhaustruct
			Name: name,
		}

		if err := mapstructure.Decode(values, &space); err != nil {
			return nil, errors.Wrap(err, "store: decode mapstructure")
		}

		spaces = append(spaces, &space)
	}

	return spaces, nil
}
