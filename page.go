package cms

import (
	"net/http"
	"strings"
	"time"

	"github.com/gosimple/slug"
)

const (
	PageRouteCMS = "page_slug"

	PageAliasPrefix    = "_page_alias_"
	PageInternalPrefix = "_page_internal_"
	PageErrorPrefix    = PageInternalPrefix + "error_"
)

type Page struct {
	ID             int64             `json:"id,omitempty"`
	SiteID         int64             `json:"site_id,omitempty"`
	ParentID       *int64            `json:"parent_id,omitempty"`
	Name           string            `json:"name,omitempty"`
	Title          string            `json:"title,omitempty"`
	RouteName      string            `json:"route_name,omitempty"`
	PageAlias      string            `json:"page_alias,omitempty"`
	Slug           string            `json:"slug,omitempty"`
	URL            string            `json:"url,omitempty"`
	CustomURL      string            `json:"custom_url,omitempty"`
	Javascript     string            `json:"javascript,omitempty"`
	Stylesheet     string            `json:"stylesheet,omitempty"`
	TemplateCode   string            `json:"template_code,omitempty"`
	Decorate       bool              `json:"decorate,omitempty"`
	Position       int               `json:"position,omitempty"`
	RequestMethods []string          `json:"request_methods,omitempty"`
	Headers        map[string]string `json:"headers,omitempty"`
	Site           *Site             `json:"site,omitempty"`
	Parent         *Page             `json:"parent,omitempty"`
	Children       []*Page           `json:"children,omitempty"`
	Metas          []Meta            `json:"metas,omitempty"`
	Metadata       map[string]any    `json:"metadata,omitempty"`
	Created        time.Time         `json:"created,omitempty"`
	Updated        time.Time         `json:"updated,omitempty"`
	Published      *time.Time        `json:"published,omitempty"`
	Expired        *time.Time        `json:"expired,omitempty"`
}

func (p *Page) String() string {
	if p.Name == "" {
		return "n/a"
	}
	return p.Name
}

func (p *Page) IsEnabled(now time.Time) bool {
	return p.Published != nil &&
		!p.Published.IsZero() &&
		(p.Published.Before(now) || p.Published.Equal(now)) &&
		(p.Expired == nil || p.Expired.IsZero() || p.Expired.After(now))
}

func (p *Page) SetAlias(pageAlias string) {
	if !strings.HasPrefix(pageAlias, PageAliasPrefix) {
		pageAlias = PageAliasPrefix + pageAlias
	}

	p.PageAlias = pageAlias
}

func (p *Page) SetInternal(routeName string) {
	if !strings.HasPrefix(routeName, PageInternalPrefix) {
		routeName = PageAliasPrefix + routeName
	}

	p.RouteName = routeName
}

func (p *Page) SetError(routeName string) {
	if !strings.HasPrefix(routeName, PageErrorPrefix) {
		routeName = PageErrorPrefix + routeName
	}
	http.NewServeMux()
	p.RouteName = routeName
}

func (p *Page) IsInternal() bool {
	return strings.HasPrefix(p.RouteName, PageInternalPrefix)
}

func (p *Page) IsError() bool {
	return strings.HasPrefix(p.RouteName, PageErrorPrefix)
}

func (p *Page) IsHybrid() bool {
	return PageRouteCMS != p.RouteName && !p.IsInternal()
}

func (p *Page) IsCMS() bool {
	return PageRouteCMS == p.RouteName && !p.IsInternal()
}

func (p *Page) IsDynamic() bool {
	return p.IsHybrid() && strings.ContainsAny(p.URL, ":{*")
}

func (p *Page) FixURL() {
	if p.IsInternal() {
		p.URL = ""
		return
	}

	if !p.IsHybrid() {
		if p.Parent == nil {
			p.Slug = ""
			p.URL = "/" + strings.TrimLeft(p.CustomURL, "/")
		} else {
			if p.Slug == "" {
				p.Slug = slug.Make(p.Name)
			}

			base := p.Parent.URL
			if base != "/" && !strings.HasSuffix(base, "/") {
				base += "/"
			}

			url := p.CustomURL
			if url == "" {
				url = p.Slug
			}

			p.URL = "/" + strings.TrimLeft(url, "/")
		}
	}

	for _, child := range p.Children {
		child.FixURL()
	}
}
