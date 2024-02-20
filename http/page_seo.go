package http

import (
	"strings"

	"github.com/gowool/cms"
)

var _ PageSEO = (*pageSEO)(nil)

type PageSEO interface {
	AddFromSite(site cms.Site) PageSEO
	AddFromPage(page cms.Page) PageSEO
	Title() string
	SetTitle(title string) PageSEO
	AddTitlePrefix(prefix string) PageSEO
	AddTitleSuffix(suffix string) PageSEO
	OriginalTitle() string
	SetSeparator(separator string) PageSEO
	Metas() map[string]map[string]string
	SetMetas(metas map[string]map[string]string) PageSEO
	AddMeta(typ, name, content string) PageSEO
	RemoveMeta(typ, name string) PageSEO
	HasMeta(typ, name string) bool
	HTMLAttributes() map[string]string
	SetHTMLAttributes(attrs map[string]string) PageSEO
	AddHTMLAttribute(name, content string) PageSEO
	RemoveHTMLAttribute(name string) PageSEO
	HasHTMLAttribute(name string) bool
	HeadAttributes() map[string]string
	SetHeadAttributes(attrs map[string]string) PageSEO
	AddHeadAttribute(name, content string) PageSEO
	RemoveHeadAttribute(name string) PageSEO
	HasHeadAttribute(name string) bool
	LinkCanonical() string
	SetLinkCanonical(link string) PageSEO
	RemoveLinkCanonical() PageSEO
	LangAlternates() map[string]string
	SetLangAlternates(langAlternates map[string]string) PageSEO
	AddLangAlternate(href, hreflang string) PageSEO
	RemoveLangAlternate(href string) PageSEO
	HasLangAlternate(href string) bool
	OEmbedLinks() map[string]string
	AddOEmbedLink(title, link string) PageSEO
}

func NewPageSEO() PageSEO {
	return &pageSEO{
		separator:      " - ",
		metas:          map[string]map[string]string{},
		htmlAttrs:      map[string]string{},
		headAttrs:      map[string]string{},
		langAlternates: map[string]string{},
		oembedLinks:    map[string]string{},
	}
}

type pageSEO struct {
	title          string
	originalTitle  string
	separator      string
	linkCanonical  string
	metas          map[string]map[string]string
	htmlAttrs      map[string]string
	headAttrs      map[string]string
	langAlternates map[string]string
	oembedLinks    map[string]string
}

func (seo *pageSEO) AddFromSite(site cms.Site) PageSEO {
	if site.Title != "" {
		seo.SetTitle(site.Title)
	}

	if site.Separator != "" {
		seo.SetSeparator(site.Separator)
	}

	if site.Locale != "" {
		lang, _, _ := strings.Cut(site.Locale, "_")
		seo.AddHTMLAttribute("lang", lang)
	}

	seo.setMetas(site.Metas)

	return seo
}

func (seo *pageSEO) AddFromPage(page cms.Page) PageSEO {
	if page.Title != "" {
		seo.AddTitlePrefix(page.Title)
	}

	seo.AddMeta(cms.MetaProperty.String(), "og:type", "website")
	seo.setMetas(page.Metas)

	seo.AddHTMLAttribute("prefix", "og: http://ogp.me/ns#")

	return seo
}

func (seo *pageSEO) setMetas(metas []cms.Meta) {
	for _, meta := range metas {
		if meta.Content == "" && (meta.Key == "keywords" || meta.Key == "description") {
			continue
		}

		seo.AddMeta(meta.Type.String(), meta.Key, meta.Content)
	}
}

func (seo *pageSEO) Title() string {
	return seo.title
}

func (seo *pageSEO) SetTitle(title string) PageSEO {
	seo.title = title
	seo.originalTitle = title
	return seo
}

func (seo *pageSEO) AddTitlePrefix(prefix string) PageSEO {
	seo.title = prefix + seo.separator + seo.title
	return seo
}

func (seo *pageSEO) AddTitleSuffix(suffix string) PageSEO {
	seo.title += seo.separator + suffix
	return seo
}

func (seo *pageSEO) OriginalTitle() string {
	return seo.originalTitle
}

func (seo *pageSEO) SetSeparator(separator string) PageSEO {
	seo.separator = separator
	return seo
}

func (seo *pageSEO) Metas() map[string]map[string]string {
	return seo.metas
}

func (seo *pageSEO) SetMetas(metas map[string]map[string]string) PageSEO {
	seo.metas = metas
	return seo
}

func (seo *pageSEO) AddMeta(typ, name, content string) PageSEO {
	if _, ok := seo.metas[typ]; !ok {
		seo.metas[typ] = map[string]string{}
	}

	seo.metas[typ][name] = content
	return seo
}

func (seo *pageSEO) RemoveMeta(typ, name string) PageSEO {
	delete(seo.metas[typ], name)
	return seo
}

func (seo *pageSEO) HasMeta(typ, name string) bool {
	if _, ok := seo.metas[typ]; !ok {
		return false
	}

	_, ok := seo.metas[typ][name]
	return ok
}

func (seo *pageSEO) HTMLAttributes() map[string]string {
	return seo.htmlAttrs
}

func (seo *pageSEO) SetHTMLAttributes(attrs map[string]string) PageSEO {
	seo.htmlAttrs = attrs
	return seo
}

func (seo *pageSEO) AddHTMLAttribute(name, content string) PageSEO {
	seo.htmlAttrs[name] = content
	return seo
}

func (seo *pageSEO) RemoveHTMLAttribute(name string) PageSEO {
	delete(seo.htmlAttrs, name)
	return seo
}

func (seo *pageSEO) HasHTMLAttribute(name string) bool {
	_, ok := seo.htmlAttrs[name]
	return ok
}

func (seo *pageSEO) HeadAttributes() map[string]string {
	return seo.headAttrs
}

func (seo *pageSEO) SetHeadAttributes(attrs map[string]string) PageSEO {
	seo.headAttrs = attrs
	return seo
}

func (seo *pageSEO) AddHeadAttribute(name, content string) PageSEO {
	seo.headAttrs[name] = content
	return seo
}

func (seo *pageSEO) RemoveHeadAttribute(name string) PageSEO {
	delete(seo.headAttrs, name)
	return seo
}

func (seo *pageSEO) HasHeadAttribute(name string) bool {
	_, ok := seo.headAttrs[name]
	return ok
}

func (seo *pageSEO) LinkCanonical() string {
	return seo.linkCanonical
}

func (seo *pageSEO) SetLinkCanonical(link string) PageSEO {
	seo.linkCanonical = link
	return seo
}

func (seo *pageSEO) RemoveLinkCanonical() PageSEO {
	seo.linkCanonical = ""
	return seo
}

func (seo *pageSEO) LangAlternates() map[string]string {
	return seo.langAlternates
}

func (seo *pageSEO) SetLangAlternates(langAlternates map[string]string) PageSEO {
	seo.langAlternates = langAlternates
	return seo
}

func (seo *pageSEO) AddLangAlternate(href, hreflang string) PageSEO {
	seo.langAlternates[href] = hreflang
	return seo
}

func (seo *pageSEO) RemoveLangAlternate(href string) PageSEO {
	delete(seo.langAlternates, href)
	return seo
}

func (seo *pageSEO) HasLangAlternate(href string) bool {
	_, ok := seo.langAlternates[href]
	return ok
}

func (seo *pageSEO) OEmbedLinks() map[string]string {
	return seo.oembedLinks
}

func (seo *pageSEO) AddOEmbedLink(title, link string) PageSEO {
	seo.oembedLinks[title] = link
	return seo
}
