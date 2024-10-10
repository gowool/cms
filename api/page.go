package api

import (
	"context"
	"maps"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/labstack/echo/v4"

	"github.com/gowool/cms"
	"github.com/gowool/cms/model"
	"github.com/gowool/cms/repository"
)

type PageBody struct {
	SiteID     int64             `json:"site_id,omitempty" yaml:"site_id,omitempty" required:"true"`
	ParentID   *int64            `json:"parent_id,omitempty" yaml:"parent_id,omitempty" required:"false"`
	Name       string            `json:"name,omitempty" yaml:"name,omitempty" required:"true"`
	Title      string            `json:"title,omitempty" yaml:"title,omitempty" required:"false"`
	Pattern    string            `json:"pattern,omitempty" yaml:"pattern,omitempty" required:"true"`
	Alias      string            `json:"alias,omitempty" yaml:"alias,omitempty" required:"false"`
	Slug       string            `json:"slug,omitempty" yaml:"slug,omitempty" required:"false"`
	CustomURL  string            `json:"custom_url,omitempty" yaml:"custom_url,omitempty" required:"false"`
	Javascript string            `json:"javascript,omitempty" yaml:"javascript,omitempty" required:"false"`
	Stylesheet string            `json:"stylesheet,omitempty" yaml:"stylesheet,omitempty" required:"false"`
	Template   string            `json:"template,omitempty" yaml:"template,omitempty" required:"true"`
	Decorate   bool              `json:"decorate,omitempty" yaml:"decorate,omitempty" required:"false"`
	Position   int               `json:"position,omitempty" yaml:"position,omitempty" required:"false"`
	Headers    map[string]string `json:"headers,omitempty" yaml:"headers,omitempty" required:"false"`
	Metas      []model.Meta      `json:"metas,omitempty" yaml:"metas,omitempty" required:"false"`
	Metadata   map[string]string `json:"metadata,omitempty" yaml:"metadata,omitempty" required:"false"`
	Published  *time.Time        `json:"published,omitempty" yaml:"published,omitempty" required:"false"`
	Expired    *time.Time        `json:"expired,omitempty" yaml:"expired,omitempty" required:"false"`
}

func (dto PageBody) Decode(m *model.Page) {
	m.SiteID = dto.SiteID
	m.ParentID = dto.ParentID
	m.Name = dto.Name
	m.Title = dto.Title
	m.Pattern = dto.Pattern
	m.Alias = dto.Alias
	m.Slug = dto.Slug
	m.CustomURL = dto.CustomURL
	m.Javascript = dto.Javascript
	m.Stylesheet = dto.Stylesheet
	m.Template = dto.Template
	m.Decorate = dto.Decorate
	m.Position = dto.Position
	m.Headers = dto.Headers
	m.Metas = dto.Metas
	m.Metadata = dto.Metadata
	m.Published = dto.Published
	m.Expired = dto.Expired
}

type Page struct {
	CRUD[PageBody, model.Page, int64]
	cfgRepo repository.Configuration
}

func NewPage(pageRepo repository.Page, cfgRepo repository.Configuration, errorTransformer ErrorTransformerFunc) Page {
	return Page{
		CRUD:    NewCRUD[PageBody](pageRepo, errorTransformer, "/pages", "Page", "Pages", "Page"),
		cfgRepo: cfgRepo,
	}
}

func (h Page) Register(e *echo.Echo, api huma.API) {
	h.CRUD.Register(e, api)

	Register(api, h.HybridPatterns(e), huma.Operation{
		Summary: "Get Hybrid Patterns",
		Method:  http.MethodGet,
		Path:    h.Path + "/hybrid-patterns",
		Tags:    h.Tags,
		Metadata: map[string]any{
			"target": cms.NewCallTarget(cms.AccessRead),
		},
	})
}

type Route struct {
	Pattern string   `json:"pattern" yaml:"pattern" required:"true"`
	Methods []string `json:"methods,omitempty" yaml:"methods,omitempty" required:"false"`
}

func (h Page) HybridPatterns(e *echo.Echo) func(context.Context, *struct{}) (*Response[[]*Route], error) {
	return func(ctx context.Context, _ *struct{}) (*Response[[]*Route], error) {
		cfg, err := h.cfgRepo.Load(ctx)
		if err != nil {
			return nil, err
		}

		routes := map[string]*Route{}
		for _, r := range e.Routes() {
			if r.Method == "echo_route_not_found" {
				continue
			}
			if cfg.IgnorePattern(r.Path) {
				continue
			}

			p, ok := routes[r.Path]
			if !ok {
				p = &Route{Pattern: r.Path}
				routes[r.Path] = p
			}
			if r.Method != "" {
				p.Methods = append(p.Methods, r.Method)
			}
		}
		return &Response[[]*Route]{
			Body: slices.SortedFunc(maps.Values(routes), func(i, j *Route) int {
				return strings.Compare(i.Pattern, j.Pattern)
			}),
		}, nil
	}
}
