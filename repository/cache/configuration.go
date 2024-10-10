package cache

import (
	"context"

	"github.com/gowool/cms"
	"github.com/gowool/cms/model"
	"github.com/gowool/cms/repository"
)

type ConfigurationRepository struct {
	repository.Configuration
	cache cms.Cache
	key   string
}

func NewConfigurationRepository(inner repository.Configuration, c cms.Cache) ConfigurationRepository {
	return ConfigurationRepository{
		Configuration: inner,
		cache:         c,
		key:           "cms::page:configuration",
	}
}

func (r ConfigurationRepository) Load(ctx context.Context) (m model.Configuration, err error) {
	if err = r.cache.Get(ctx, r.key, &m); err == nil {
		return
	}

	if m, err = r.Configuration.Load(ctx); err != nil {
		return
	}

	_ = r.cache.Set(ctx, r.key, m)
	return
}

func (r ConfigurationRepository) Save(ctx context.Context, m *model.Configuration) error {
	defer func() {
		_ = r.cache.DelByKey(ctx, r.key)
	}()

	return r.Configuration.Save(ctx, m)
}
