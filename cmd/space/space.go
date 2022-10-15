// Package space contains wrapper for adding/removing values from Spaces.
package space

import (
	"io/fs"

	"github.com/cockroachdb/errors"
	"github.com/shihanng/tfspace/config"
	"github.com/shihanng/tfspace/space"
	"github.com/shihanng/tfspace/store"
	"go.uber.org/zap"
)

// WithSpace wraps around exec. If load tfspace.yml then execute exec,
// then save the changes back to tfspace.yml.
func WithSpace(exec func(s *space.Spaces) error) error {
	logger := zap.L()

	cfg, err := config.GetConfig()
	if err != nil {
		return err
	}

	logger = logger.With(zap.String("config_path", cfg.Path))

	logger.Debug("Load spaces")

	spaces, err := store.Load(cfg.Path)
	if err != nil {
		if !errors.Is(err, fs.ErrNotExist) {
			return errors.Wrap(err, "space: fail to load spaces")
		}

		logger.Debug("Config does not exist")

		spaces = space.Spaces{}
	}

	if err := exec(&spaces); err != nil {
		return err
	}

	if err := store.Save(cfg.Path, spaces); err != nil {
		return err
	}

	return nil
}
