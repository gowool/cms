package http

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"slices"
	"sort"
	"strings"

	"github.com/dlclark/regexp2"

	"github.com/gowool/cms"
	"github.com/gowool/cms/internal"
)

var (
	_ SiteSelector = HostSiteSelector{}
	_ SiteSelector = HostByLocaleSiteSelector{}
	_ SiteSelector = HostPathSiteSelector{}
	_ SiteSelector = HostPathByLocaleSiteSelector{}
)

const (
	reOptions    = regexp2.IgnoreCase & regexp2.RE2
	rePathExpr   = "^(%s)(/.*|$)"
	reNoPathExpr = "^()(/.*|$)"
)

var ErrSiteNotFound = errors.New("site not found")

var reNoPath *regexp2.Regexp

func init() {
	var err error
	reNoPath, err = regexp2.Compile(reNoPathExpr, reOptions)
	if err != nil {
		panic(err)
	}
}

type SelectorType string

func (t SelectorType) String() string {
	return string(t)
}

const (
	HostSelector             SelectorType = "host"
	HostByLocaleSelector     SelectorType = "host_by_locale"
	HostPathSelector         SelectorType = "host_with_path"
	HostPathByLocaleSelector SelectorType = "host_with_path_by_locale"
)

var Selectors = []SelectorType{
	HostSelector,
	HostByLocaleSelector,
	HostPathSelector,
	HostPathByLocaleSelector,
}

type SiteSelector interface {
	Type() SelectorType
	Select(r *http.Request) (cms.Site, string, error)
}

type HostSiteSelector struct {
	siteService cms.SiteService
}

func NewHostSiteSelector(siteService cms.SiteService) HostSiteSelector {
	return HostSiteSelector{
		siteService: siteService,
	}
}

func (s HostSiteSelector) Type() SelectorType {
	return HostSelector
}

func (s HostSiteSelector) Select(r *http.Request) (cms.Site, string, error) {
	host := Host(r)
	sites, err := s.siteService.GetByHost(r.Context(), host)
	if err != nil {
		return cms.Site{}, "", err
	}

	var site cms.Site
	for _, site = range sites {
		if !site.IsLocalhost() {
			break
		}
	}
	return site, "", nil
}

type HostByLocaleSiteSelector struct {
	siteService cms.SiteService
}

func NewHostByLocaleSiteSelector(siteService cms.SiteService) HostByLocaleSiteSelector {
	return HostByLocaleSiteSelector{
		siteService: siteService,
	}
}

func (s HostByLocaleSiteSelector) Type() SelectorType {
	return HostByLocaleSelector
}

func (s HostByLocaleSiteSelector) Select(r *http.Request) (cms.Site, string, error) {
	host := Host(r)
	sites, err := s.siteService.GetByHost(r.Context(), host)
	if err != nil {
		return cms.Site{}, "", err
	}

	preferredSites := make([]cms.Site, 0, len(sites))
	for _, site := range sites {
		preferredSites = append(preferredSites, site)
		if !site.IsLocalhost() {
			break
		}
	}
	slices.Clip(preferredSites)

	if site := PreferredSite(preferredSites, r); site.ID != 0 {
		return site, "", nil
	}

	return cms.Site{}, "", ErrSiteNotFound
}

type HostPathSiteSelector struct {
	siteService cms.SiteService
}

func NewHostPathSiteSelector(siteService cms.SiteService) HostPathSiteSelector {
	return HostPathSiteSelector{
		siteService: siteService,
	}
}

func (s HostPathSiteSelector) Type() SelectorType {
	return HostPathSelector
}

func (s HostPathSiteSelector) Select(r *http.Request) (cms.Site, string, error) {
	host := Host(r)
	sites, err := s.siteService.GetByHost(r.Context(), host)
	if err != nil {
		return cms.Site{}, "", err
	}

	var (
		site        cms.Site
		defaultSite cms.Site
	)
	path := "/"

	for _, _site := range sites {
		if defaultSite.ID == 0 && _site.IsDefault {
			defaultSite = _site
		}

		match, err := MatchRequest(_site, r)
		if err != nil {
			continue
		}

		site = _site
		path = match

		if !site.IsLocalhost() {
			break
		}
	}

	if site.ID == 0 {
		if defaultSite.ID != 0 {
			return cms.Site{}, getURL(r, defaultSite, host), nil
		}
		return cms.Site{}, "", ErrSiteNotFound
	}

	replacePath(r, path)

	return site, "", nil
}

type HostPathByLocaleSiteSelector struct {
	siteService cms.SiteService
}

func NewHostPathByLocaleSiteSelector(siteService cms.SiteService) HostPathByLocaleSiteSelector {
	return HostPathByLocaleSiteSelector{
		siteService: siteService,
	}
}

func (s HostPathByLocaleSiteSelector) Type() SelectorType {
	return HostPathByLocaleSelector
}

func (s HostPathByLocaleSiteSelector) Select(r *http.Request) (cms.Site, string, error) {
	host := Host(r)
	sites, err := s.siteService.GetByHost(r.Context(), host)
	if err != nil {
		return cms.Site{}, "", err
	}

	var site cms.Site

	path := "/"
	preferredSites := make([]cms.Site, 0, len(sites))
	for _, _site := range sites {
		preferredSites = append(preferredSites, _site)

		match, err := MatchRequest(_site, r)
		if err != nil {
			continue
		}

		site = _site
		path = match

		if !site.IsLocalhost() {
			break
		}
	}
	slices.Clip(preferredSites)

	if site.ID == 0 {
		if defaultSite := PreferredSite(preferredSites, r); defaultSite.ID != 0 {
			return cms.Site{}, getURL(r, defaultSite, host), nil
		}
		return cms.Site{}, "", ErrSiteNotFound
	}

	replacePath(r, path)

	return site, "", nil
}

func PreferredSite(sites []cms.Site, r *http.Request) cms.Site {
	if len(sites) == 0 {
		return cms.Site{}
	}

	var locales []string
	for _, site := range sites {
		if site.Locale != "" {
			locales = append(locales, site.Locale)
		}
	}

	language := PreferredLanguage(locales, r)
	host := Host(r)

	if index := slices.IndexFunc(sites, func(item cms.Site) bool {
		return item.Locale == language && (item.Host == host || item.Host == "localhost")
	}); index >= 0 {
		return sites[index]
	}
	return cms.Site{}
}

func MatchRequest(site cms.Site, r *http.Request) (string, error) {
	var (
		re    *regexp2.Regexp
		match *regexp2.Match
		err   error
	)
	if site.RelativePath == "" || site.RelativePath == "/" {
		re = reNoPath
	} else {
		re, err = regexp2.Compile(fmt.Sprintf(rePathExpr, site.RelativePath), reOptions)
		if err != nil {
			return "", err
		}
	}

	match, err = re.FindStringMatch(r.URL.Path)
	if err != nil {
		return "", err
	}

	if match == nil {
		return "", fmt.Errorf("invalid path %s", r.URL.Path)
	}

	groups := match.Groups()

	if len(groups) < 3 {
		return "", fmt.Errorf("invalid match path %s", r.URL.Path)
	}

	return groups[2].String(), nil
}

func PreferredLanguage(locales []string, r *http.Request) string {
	languages := Languages(r)

	if len(locales) == 0 {
		if len(languages) > 0 {
			return languages[0]
		}
		return ""
	}

	if len(languages) == 0 {
		return locales[0]
	}

	mapLocales := make(map[string]struct{})
	for _, locale := range locales {
		mapLocales[locale] = struct{}{}
	}

	for _, language := range languages {
		if _, ok := mapLocales[language]; ok {
			return language
		}

		if codes := strings.Split(language, "_"); len(codes) > 1 {
			if _, ok := mapLocales[codes[0]]; ok {
				return codes[0]
			}
		}
	}

	return ""
}

type lang struct {
	l string
	q string
}

func Languages(r *http.Request) []string {
	header := r.Header.Get("Accept-Language")

	languages := internal.Map(strings.Split(header, ","), func(item string) lang {
		data := strings.Split(strings.TrimSpace(item), ";q=")
		codes := strings.Split(data[0], "-")

		if codes[0][0] == 'i' {
			// Language not listed in ISO 639 that are not variants
			// of any listed language, which can be registered with the
			// i-prefix, such as i-cherokee
			if len(codes[0]) > 1 {
				codes[0] = codes[0][1:]
			}
		}

		if len(codes) > 1 {
			codes[1] = strings.ToUpper(codes[1])
		}

		l := lang{l: strings.Join(codes, "_"), q: "1.0"}

		if len(data) > 1 {
			l.q = data[1]
			return l
		}

		if data[0] == "*" {
			l.q = "0.0"
			return l
		}

		return l
	})

	sort.Slice(languages, func(i, j int) bool {
		return languages[i].q > languages[j].q
	})

	return internal.Map(languages, func(item lang) string {
		return item.l
	})
}

func Host(r *http.Request) string {
	host, port, err := net.SplitHostPort(r.Host)
	if err != nil {
		panic(err)
	}
	if port == "80" || port == "443" {
		return host
	}
	return host + ":" + port
}

func getURL(r *http.Request, site cms.Site, host string) string {
	if site.IsLocalhost() {
		return fmt.Sprintf("%s://%s%s", Scheme(r), host, site.RelativePath)
	}
	return fmt.Sprintf("%s:%s", Scheme(r), site.URL())
}

func replacePath(r *http.Request, path string) {
	uri := path
	if r.URL.RawQuery != "" {
		uri += "?" + r.URL.RawQuery
	}

	r.RequestURI = uri
	r.URL.Path = path
	r.URL.RawPath = ""
}
