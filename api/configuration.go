package api

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/labstack/echo/v4"

	"github.com/gowool/cms"
	"github.com/gowool/cms/model"
	"github.com/gowool/cms/repository"
)

type Configuration struct {
	errorTransformer func(context.Context, error) error
	repo             repository.Configuration
	path             string
	tags             []string
}

func NewConfiguration(repo repository.Configuration, errorTransformer func(context.Context, error) error) Configuration {
	return Configuration{
		errorTransformer: errorTransformer,
		repo:             repo,
		path:             "/pages/configuration",
		tags:             []string{"Pages"},
	}
}

func (h Configuration) Register(_ *echo.Echo, api huma.API) {
	Register(api, h.load, huma.Operation{
		Summary: "Get Configuration",
		Method:  http.MethodGet,
		Path:    h.path,
		Tags:    h.tags,
		Metadata: map[string]any{
			"target": cms.NewCallTarget(cms.AccessRead),
		},
	})
	Register(api, h.save, huma.Operation{
		Summary: "Save Configuration",
		Method:  http.MethodPatch,
		Path:    h.path,
		Tags:    h.tags,
		Metadata: map[string]any{
			"target": cms.NewCallTarget(cms.AccessWrite),
		},
	})
}

func (h Configuration) load(ctx context.Context, _ *struct{}) (*Response[model.Configuration], error) {
	cfg, err := h.repo.Load(ctx)
	if err != nil {
		return nil, h.errorTransformer(ctx, err)
	}
	return &Response[model.Configuration]{Body: cfg}, nil
}

func (h Configuration) save(ctx context.Context, in *CreateInput[model.Configuration]) (*struct{}, error) {
	cfg, err := h.repo.Load(ctx)
	if err != nil {
		return nil, h.errorTransformer(ctx, err)
	}

	cfg = cfg.With(in.Body)
	if err = h.repo.Save(ctx, &cfg); err != nil {
		return nil, h.errorTransformer(ctx, err)
	}
	return nil, nil
}
