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

type Admin struct {
	List[model.Admin]
	Read[model.Admin, int64]

	path   string
	pathID string
	tags   []string
}

func NewAdmin(r repository.Admin, errorTransformer ErrorTransformerFunc) Admin {
	return Admin{
		List:   NewList(r.FindAndCount, errorTransformer),
		Read:   NewRead(r.FindByID, errorTransformer),
		path:   "/admin",
		pathID: "/admin/{id}",
		tags:   []string{"Admin"},
	}
}

func (r Admin) Register(_ *echo.Echo, humaAPI huma.API) {
	Register(humaAPI, r.List.Handler, huma.Operation{
		Summary: "Get Admins",
		Method:  http.MethodGet,
		Path:    r.path,
		Tags:    r.tags,
		Metadata: map[string]any{
			"target": cms.NewCallTarget(cms.AccessAdmin),
		},
	})
	Register(humaAPI, r.Read.Handler, huma.Operation{
		Summary: "Get Admin",
		Method:  http.MethodGet,
		Path:    r.pathID,
		Tags:    r.tags,
		Metadata: map[string]any{
			"target": cms.NewCallTarget(cms.AccessAdmin),
		},
	})
	Register(humaAPI, r.me, huma.Operation{
		Summary: "Me",
		Method:  http.MethodGet,
		Path:    r.path + "/me",
		Tags:    r.tags,
		Metadata: map[string]any{
			"target": cms.NewCallTarget(cms.AccessRead),
		},
	})
}

func (r Admin) me(ctx context.Context, _ *struct{}) (*Response[*model.Admin], error) {
	admin := cms.CtxAdmin(ctx)

	return &Response[*model.Admin]{Body: admin}, nil
}
