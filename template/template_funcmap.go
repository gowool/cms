package template

import (
	"fmt"
	"html"
	"html/template"
	"strings"

	"github.com/gowool/cms/http"
)

type funcMap struct{}

func FuncMap() template.FuncMap {
	return funcMap{}.FuncMap()
}

func (fm funcMap) FuncMap() template.FuncMap {
	return template.FuncMap{
		"title_tag":       fm.titleTag,
		"meta_tags":       fm.metaTags,
		"html_attrs":      fm.htmlAttrs,
		"head_attrs":      fm.headAttrs,
		"link_canonical":  fm.linkCanonical,
		"lang_alternates": fm.langAlternates,
		"oembed_links":    fm.oEmbedLinks,
	}
}

func (fm funcMap) titleTag(seo http.PageSEO) template.HTML {
	return template.HTML(
		fmt.Sprintf(
			"<title>%s</title>",
			StripTags(seo.Title()),
		),
	)
}

func (fm funcMap) metaTags(seo http.PageSEO) template.HTML {
	normalize := func(s string) string {
		return EscapeDoubleQuotes(StripTags(s))
	}

	var b strings.Builder
	for typ, metas := range seo.Metas() {
		for name, content := range metas {
			b.WriteString("<meta ")
			b.WriteString(typ)
			b.WriteString(`="`)
			b.WriteString(normalize(name))
			if content != "" {
				b.WriteString(`" content="`)
				b.WriteString(normalize(content))
			}
			b.WriteString("\" />\n")
		}
	}
	return template.HTML(b.String())
}

func (fm funcMap) htmlAttrs(seo http.PageSEO) string {
	return fm.attrs(seo.HTMLAttributes())
}

func (fm funcMap) headAttrs(seo http.PageSEO) string {
	return fm.attrs(seo.HeadAttributes())
}

func (fm funcMap) attrs(attrs map[string]string) string {
	var b strings.Builder
	for name, value := range attrs {
		b.WriteString(name)
		b.WriteString(`="`)
		b.WriteString(html.EscapeString(value))
		b.WriteString(`" `)
	}
	return strings.TrimRight(b.String(), " ")
}

func (fm funcMap) linkCanonical(seo http.PageSEO) template.HTML {
	if seo.LinkCanonical() != "" {
		return template.HTML(fmt.Sprintf(`<link rel="canonical" href="%s" />`, html.EscapeString(seo.LinkCanonical())))
	}
	return ""
}

func (fm funcMap) langAlternates(seo http.PageSEO) template.HTML {
	var b strings.Builder
	for href, hreflang := range seo.LangAlternates() {
		b.WriteString(`<link rel="alternate" href="`)
		b.WriteString(html.EscapeString(href))
		b.WriteString(`" hreflang="`)
		b.WriteString(html.EscapeString(hreflang))
		b.WriteString("\" />\n")
	}
	return template.HTML(b.String())
}

func (fm funcMap) oEmbedLinks(seo http.PageSEO) template.HTML {
	var b strings.Builder
	for title, link := range seo.OEmbedLinks() {
		b.WriteString(`<link rel="alternate" type="application/json+oembed" href="`)
		b.WriteString(link)
		b.WriteString(`" title="`)
		b.WriteString(html.EscapeString(title))
		b.WriteString("\" />\n")
	}
	return template.HTML(b.String())
}
