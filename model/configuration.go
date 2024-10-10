package model

import (
	"maps"
	"net/http"
	"slices"

	"github.com/gowool/cms/internal"
)

var MultisiteStrategies = []MultisiteStrategy{Host, HostByLocale, HostWithPath, HostWithPathByLocale}

const (
	Host                 = MultisiteStrategy("host")
	HostByLocale         = MultisiteStrategy("host_by_locale")
	HostWithPath         = MultisiteStrategy("host_with_path")
	HostWithPathByLocale = MultisiteStrategy("host_with_path_by_locale")
)

type MultisiteStrategy string

func (t MultisiteStrategy) IsZero() bool {
	return t == ""
}

func (t MultisiteStrategy) String() string {
	return string(t)
}

type Configuration struct {
	Debug                 bool              `json:"debug,omitempty" yaml:"debug,omitempty" required:"true"`
	Multisite             MultisiteStrategy `json:"multisite,omitempty" yaml:"multisite,omitempty" required:"false" enum:"host,host_by_locale,host_with_path,host_with_path_by_locale"`
	FallbackLocale        string            `json:"fallback_locale,omitempty" yaml:"fallback_locale,omitempty" required:"false"`
	IgnoreRequestPatterns []string          `json:"ignore_request_patterns,omitempty" yaml:"ignore_request_patterns,omitempty" required:"false"`
	IgnoreRequestURIs     []string          `json:"ignore_request_uris,omitempty" yaml:"ignore_request_uris,omitempty" required:"false"`
	CatchErrors           map[string][]int  `json:"catch_errors,omitempty" yaml:"catch_errors,omitempty" required:"false"`
	Additional            map[string]string `json:"additional,omitempty" yaml:"additional,omitempty" required:"false"`
}

func NewConfiguration() Configuration {
	return Configuration{
		Multisite:      Host,
		FallbackLocale: "en_US",
		CatchErrors: map[string][]int{
			PageError4xx: {
				http.StatusBadRequest,
				http.StatusUnauthorized,
				http.StatusPaymentRequired,
				http.StatusForbidden,
				http.StatusNotFound,
				http.StatusMethodNotAllowed,
				http.StatusNotAcceptable,
				http.StatusProxyAuthRequired,
				http.StatusRequestTimeout,
				http.StatusConflict,
				http.StatusGone,
				http.StatusLengthRequired,
				http.StatusPreconditionFailed,
				http.StatusRequestEntityTooLarge,
				http.StatusRequestURITooLong,
				http.StatusUnsupportedMediaType,
				http.StatusRequestedRangeNotSatisfiable,
				http.StatusExpectationFailed,
				http.StatusTeapot,
				http.StatusMisdirectedRequest,
				http.StatusUnprocessableEntity,
				http.StatusLocked,
				http.StatusFailedDependency,
				http.StatusTooEarly,
				http.StatusUpgradeRequired,
				http.StatusPreconditionRequired,
				http.StatusTooManyRequests,
				http.StatusRequestHeaderFieldsTooLarge,
				http.StatusUnavailableForLegalReasons,
			},
			PageError5xx: {
				http.StatusInternalServerError,
				http.StatusNotImplemented,
				http.StatusBadGateway,
				http.StatusServiceUnavailable,
				http.StatusGatewayTimeout,
				http.StatusHTTPVersionNotSupported,
				http.StatusVariantAlsoNegotiates,
				http.StatusInsufficientStorage,
				http.StatusLoopDetected,
				http.StatusNotExtended,
				http.StatusNetworkAuthenticationRequired,
			},
		},
		Additional: map[string]string{},
	}
}

func (c Configuration) IgnorePattern(pattern string) bool {
	if pattern == "" {
		return false
	}

	for _, expr := range c.IgnoreRequestPatterns {
		if re, ok := internal.Regexp(expr); ok {
			if ok, _ = re.MatchString(pattern); ok {
				return true
			}
		}
	}
	return false
}

func (c Configuration) IgnoreURI(uri string) bool {
	for _, expr := range c.IgnoreRequestURIs {
		if re, ok := internal.Regexp(expr); ok {
			if ok, _ = re.MatchString(uri); ok {
				return true
			}
		}
	}
	return false
}

func (c Configuration) With(other Configuration) Configuration {
	c.Debug = other.Debug

	if !other.Multisite.IsZero() {
		c.Multisite = other.Multisite
	}

	c.IgnoreRequestPatterns = internal.Unique(slices.Concat(c.IgnoreRequestPatterns, other.IgnoreRequestPatterns))
	c.IgnoreRequestURIs = internal.Unique(slices.Concat(c.IgnoreRequestURIs, other.IgnoreRequestURIs))

	if other.FallbackLocale != "" {
		c.FallbackLocale = other.FallbackLocale
	}

	if other.CatchErrors != nil {
		if c.CatchErrors == nil {
			c.CatchErrors = make(map[string][]int)
		}
		for pattern, codes := range other.CatchErrors {
			c.CatchErrors[pattern] = internal.Unique(slices.Concat(c.CatchErrors[pattern], codes))
		}
	}

	if other.Additional != nil {
		if c.Additional == nil {
			c.Additional = make(map[string]string)
		}
		maps.Copy(c.Additional, other.Additional)
	}
	return c
}
