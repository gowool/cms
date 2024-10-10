package fallback

import (
	"context"

	"github.com/gowool/cms/model"
	"github.com/gowool/cms/repository"
)

type ConfigurationRepository struct {
	repository.Configuration
	cfg model.Configuration
}

func NewConfigurationRepository(inner repository.Configuration, cfg model.Configuration) ConfigurationRepository {
	return ConfigurationRepository{
		Configuration: inner,
		cfg:           cfg,
	}
}

func (r ConfigurationRepository) Load(ctx context.Context) (model.Configuration, error) {
	if m, err := r.Configuration.Load(ctx); err == nil {
		return r.cfg.With(m), nil
	}
	return r.cfg, nil
}
