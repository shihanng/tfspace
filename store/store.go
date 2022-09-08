package store

import (
	"os"

	"github.com/cockroachdb/errors"
	"github.com/goccy/go-yaml"
	"github.com/mitchellh/mapstructure"
	"github.com/shihanng/tfspace/space"
)

func Load(path string) (space.Spaces, error) {
	dat, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var v yaml.MapSlice
	if err := yaml.Unmarshal(dat, &v); err != nil {
		return nil, errors.Wrap(err, "store: decode yaml")
	}

	return spacesFromMapSlice(v)
}

func Save(path string, spaces space.Spaces) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	payload := make(yaml.MapSlice, 0, len(spaces))

	for _, space := range spaces {
		payload = append(payload, yaml.MapItem{
			Key: space.Name,
			Value: struct {
				Backend   []string
				Varfile   []string
				Workspace string
			}{
				Backend:   space.Backend,
				Varfile:   space.Varfile,
				Workspace: space.Workspace,
			},
		})
	}

	return yaml.NewEncoder(file, yaml.IndentSequence(true)).Encode(payload)
}

func spacesFromMapSlice(ms yaml.MapSlice) (space.Spaces, error) {
	spaces := make(space.Spaces, 0, len(ms))

	for _, item := range ms {
		name, ok := item.Key.(string)
		if !ok {
			return nil, errors.New("store: name is not string")
		}

		values, ok := item.Value.(map[string]interface{})
		if !ok {
			return nil, errors.New("store: value is not hash map")
		}

		sp := space.Space{
			Name: name,
		}

		if err := mapstructure.Decode(values, &sp); err != nil {
			return nil, errors.Wrap(err, "store: decode mapstructure")
		}

		spaces = append(spaces, &sp)
	}

	return spaces, nil
}
