package space

import (
	"errors"
	"io/fs"

	"github.com/shihanng/tfspace/config"
	"github.com/shihanng/tfspace/space"
	"github.com/shihanng/tfspace/store"
	"go.uber.org/zap"
)

func WithSpace(exec func(s space.Spaces)) error {
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
			return err
		}

		logger.Debug("Config does not exist")
		spaces = space.Spaces{}
	}

	exec(spaces)

	if err := store.Save(cfg.Path, spaces); err != nil {
		return err
	}

	return nil
}
