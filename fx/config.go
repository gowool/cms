package fx

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/fx"

	"github.com/gowool/cms"
)

type SameSiteType string

const (
	SameSiteDefault SameSiteType = "default"
	SameSiteLax     SameSiteType = "lax"
	SameSiteStrict  SameSiteType = "strict"
	SameSiteNone    SameSiteType = "none"
)

func (s SameSiteType) HTTP() http.SameSite {
	switch s {
	case SameSiteDefault:
		return http.SameSiteDefaultMode
	case SameSiteLax:
		return http.SameSiteLaxMode
	case SameSiteStrict:
		return http.SameSiteStrictMode
	case SameSiteNone:
		return http.SameSiteNoneMode
	default:
		panic("invalid same site")
	}
}

type AreaConfig struct {
	Enabled     bool           `json:"enabled,omitempty" yaml:"enabled,omitempty"`
	BasePath    string         `json:"base_path,omitempty" yaml:"base_path,omitempty"`
	Middlewares []string       `json:"middlewares,omitempty" yaml:"middlewares,omitempty"`
	Additional  map[string]any `json:"additional,omitempty" yaml:"additional,omitempty"`
}

type HumaAreaConfig struct {
	AreaConfig      `json:",inline" yaml:",inline"`
	HumaMiddlewares []string `json:"huma_middlewares,omitempty" yaml:"huma_middlewares,omitempty"`
}

type AreasConfig struct {
	API      HumaAreaConfig `json:"api,omitempty" yaml:"api,omitempty"`
	AdminAPI HumaAreaConfig `json:"admin_api,omitempty" yaml:"admin_api,omitempty"`
	Static   AreaConfig     `json:"static,omitempty" yaml:"static,omitempty"`
	Web      AreaConfig     `json:"web,omitempty" yaml:"web,omitempty"`
}

type SessionConfig struct {
	CleanupInterval time.Duration `json:"cleanup_interval" yaml:"cleanup_interval"`
	IdleTimeout     time.Duration `json:"idle_timeout" yaml:"idle_timeout"`
	Lifetime        time.Duration `json:"lifetime" yaml:"lifetime"`
	Cookie          struct {
		Name     string       `json:"name" yaml:"name"`
		Domain   string       `json:"domain" yaml:"domain"`
		Path     string       `json:"path" yaml:"path"`
		Persist  bool         `json:"persist" yaml:"persist"`
		Secure   bool         `json:"secure" yaml:"secure"`
		HTTPOnly bool         `json:"http_only" yaml:"http_only"`
		SameSite SameSiteType `json:"same_site" yaml:"same_site"`
	} `json:"cookie" yaml:"cookie"`
}

type MiddlewareConfig struct {
	fx.Out
	Global    GlobalConfig    `json:"global,omitempty" yaml:"global,omitempty"`
	Recover   RecoverConfig   `json:"recover,omitempty" yaml:"recover,omitempty"`
	BodyLimit BodyLimitConfig `json:"body_limit,omitempty" yaml:"body_limit,omitempty"`
	Compress  GzipConfig      `json:"compress,omitempty" yaml:"compress,omitempty"`
	Secure    SecureConfig    `json:"secure,omitempty" yaml:"secure,omitempty"`
	CORS      CORSConfig      `json:"cors,omitempty" yaml:"cors,omitempty"`
	CSRF      CSRFConfig      `json:"csrf,omitempty" yaml:"csrf,omitempty"`
}

func (cfg *MiddlewareConfig) InitDefaults() {
	cfg.BodyLimit.InitDefaults()
	cfg.Compress.InitDefaults()
	cfg.Secure.InitDefaults()
}

type GlobalConfig struct {
	BeforeRouter []string `json:"before_router,omitempty" yaml:"before_router,omitempty"`
	AfterRouter  []string `json:"after_router,omitempty" yaml:"after_router,omitempty"`
}

type RecoverConfig struct {
	// Skipper defines a function to skip middleware.
	Skipper middleware.Skipper

	// Size of the stack to be printed.
	// Optional. Default value 4KB.
	StackSize int `json:"stack_size,omitempty" yaml:"stack_size,omitempty"`

	// DisableStackAll disables formatting stack traces of all other goroutines
	// into buffer after the trace for the current goroutine.
	// Optional. Default value false.
	DisableStackAll bool `json:"disable_stack_all,omitempty" yaml:"disable_stack_all,omitempty"`
}

type BodyLimitConfig struct {
	// Skipper defines a function to skip middleware.
	Skipper middleware.Skipper `json:"-" yaml:"-"`

	// Maximum allowed size for a request body, it can be specified
	// as `4x` or `4xB`, where x is one of the multiple from K, M, G, T or P.
	Limit string `json:"limit,omitempty" yaml:"limit,omitempty"`
}

func (cfg *BodyLimitConfig) InitDefaults() {
	if cfg.Limit == "" {
		cfg.Limit = "4KB"
	}
}

type GzipConfig struct {
	// Skipper defines a function to skip middleware.
	Skipper middleware.Skipper `json:"-" yaml:"-"`

	// Gzip compression level.
	// Optional. Default value -1.
	Level int `json:"level,omitempty" yaml:"level,omitempty"`

	// Length threshold before gzip compression is applied.
	// Optional. Default value 0.
	//
	// Most of the time you will not need to change the default. Compressing
	// a short response might increase the transmitted data because of the
	// gzip format overhead. Compressing the response will also consume CPU
	// and time on the server and the client (for decompressing). Depending on
	// your use case such a threshold might be useful.
	//
	// See also:
	// https://webmasters.stackexchange.com/questions/31750/what-is-recommended-minimum-object-size-for-gzip-performance-benefits
	MinLength int `json:"min_length,omitempty" yaml:"min_length,omitempty"`
}

func (cfg *GzipConfig) InitDefaults() {
	if cfg.MinLength <= 0 {
		cfg.MinLength = 1024
	}
}

type SecureConfig struct {
	// Skipper defines a function to skip middleware.
	Skipper middleware.Skipper `json:"-" yaml:"-"`

	// XSSProtection provides protection against cross-site scripting attack (XSS)
	// by setting the `X-XSS-Protection` header.
	// Optional. Default value "1; mode=block".
	XSSProtection string `json:"xss_protection,omitempty" yaml:"xss_protection,omitempty"`

	// ContentTypeNosniff provides protection against overriding Content-Type
	// header by setting the `X-Content-Type-Options` header.
	// Optional. Default value "nosniff".
	ContentTypeNosniff string `json:"content_type_nosniff,omitempty" yaml:"content_type_nosniff,omitempty"`

	// XFrameOptions can be used to indicate whether or not a browser should
	// be allowed to render a page in a <frame>, <iframe> or <object> .
	// Sites can use this to avoid clickjacking attacks, by ensuring that their
	// content is not embedded into other sites.provides protection against
	// clickjacking.
	// Optional. Default value "SAMEORIGIN".
	// Possible values:
	// - "SAMEORIGIN" - The page can only be displayed in a frame on the same origin as the page itself.
	// - "DENY" - The page cannot be displayed in a frame, regardless of the site attempting to do so.
	// - "ALLOW-FROM uri" - The page can only be displayed in a frame on the specified origin.
	XFrameOptions string `json:"x_frame_options,omitempty" yaml:"x_frame_options,omitempty"`

	// HSTSMaxAge sets the `Strict-Transport-Security` header to indicate how
	// long (in seconds) browsers should remember that this site is only to
	// be accessed using HTTPS. This reduces your exposure to some SSL-stripping
	// man-in-the-middle (MITM) attacks.
	// Optional. Default value 0.
	HSTSMaxAge int `json:"hsts_max_age,omitempty" yaml:"hsts_max_age,omitempty"`

	// HSTSExcludeSubdomains won't include subdomains tag in the `Strict Transport Security`
	// header, excluding all subdomains from security policy. It has no effect
	// unless HSTSMaxAge is set to a non-zero value.
	// Optional. Default value false.
	HSTSExcludeSubdomains bool `json:"hsts_exclude_subdomains,omitempty" yaml:"hsts_exclude_subdomains,omitempty"`

	// ContentSecurityPolicy sets the `Content-Security-Policy` header providing
	// security against cross-site scripting (XSS), clickjacking and other code
	// injection attacks resulting from execution of malicious content in the
	// trusted web page context.
	// Optional. Default value "".
	ContentSecurityPolicy string `json:"content_security_policy,omitempty" yaml:"content_security_policy,omitempty"`

	// CSPReportOnly would use the `Content-Security-Policy-Report-Only` header instead
	// of the `Content-Security-Policy` header. This allows iterative updates of the
	// content security policy by only reporting the violations that would
	// have occurred instead of blocking the resource.
	// Optional. Default value false.
	CSPReportOnly bool `json:"csp_report_only,omitempty" yaml:"csp_report_only,omitempty"`

	// HSTSPreloadEnabled will add the preload tag in the `Strict Transport Security`
	// header, which enables the domain to be included in the HSTS preload list
	// maintained by Chrome (and used by Firefox and Safari): https://hstspreload.org/
	// Optional.  Default value false.
	HSTSPreloadEnabled bool `json:"hsts_preload_enabled,omitempty" yaml:"hsts_preload_enabled,omitempty"`

	// ReferrerPolicy sets the `Referrer-Policy` header providing security against
	// leaking potentially sensitive request paths to third parties.
	// Optional. Default value "".
	ReferrerPolicy string `json:"referrer_policy,omitempty" yaml:"referrer_policy,omitempty"`
}

func (cfg *SecureConfig) InitDefaults() {
	if cfg.XSSProtection == "" {
		cfg.XSSProtection = middleware.DefaultSecureConfig.XSSProtection
	}
	if cfg.ContentTypeNosniff == "" {
		cfg.ContentTypeNosniff = middleware.DefaultSecureConfig.ContentTypeNosniff
	}
	if cfg.XFrameOptions == "" {
		cfg.XFrameOptions = middleware.DefaultSecureConfig.XFrameOptions
	}
}

type CORSConfig struct {
	// Skipper defines a function to skip middleware.
	Skipper middleware.Skipper `json:"-" yaml:"-"`

	// AllowOrigins determines the value of the Access-Control-Allow-Origin
	// response header.  This header defines a list of origins that may access the
	// resource.  The wildcard characters '*' and '?' are supported and are
	// converted to regex fragments '.*' and '.' accordingly.
	//
	// Security: use extreme caution when handling the origin, and carefully
	// validate any logic. Remember that attackers may register hostile domain names.
	// See https://blog.portswigger.net/2016/10/exploiting-cors-misconfigurations-for.html
	//
	// Optional. Default value []string{"*"}.
	//
	// See also: https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Allow-Origin
	AllowOrigins []string `json:"allow_origins,omitempty" yaml:"allow_origins,omitempty"`

	// AllowOriginFunc is a custom function to validate the origin. It takes the
	// origin as an argument and returns true if allowed or false otherwise. If
	// an error is returned, it is returned by the handler. If this option is
	// set, AllowOrigins is ignored.
	//
	// Security: use extreme caution when handling the origin, and carefully
	// validate any logic. Remember that attackers may register hostile domain names.
	// See https://blog.portswigger.net/2016/10/exploiting-cors-misconfigurations-for.html
	//
	// Optional.
	AllowOriginFunc func(origin string) (bool, error) `json:"-" yaml:"-"`

	// AllowMethods determines the value of the Access-Control-Allow-Methods
	// response header.  This header specified the list of methods allowed when
	// accessing the resource.  This is used in response to a preflight request.
	//
	// Optional. Default value DefaultCORSConfig.AllowMethods.
	// If `allowMethods` is left empty, this middleware will fill for preflight
	// request `Access-Control-Allow-Methods` header value
	// from `Allow` header that echo.Router set into context.
	//
	// See also: https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Allow-Methods
	AllowMethods []string `json:"allow_methods,omitempty" yaml:"allow_methods,omitempty"`

	// AllowHeaders determines the value of the Access-Control-Allow-Headers
	// response header.  This header is used in response to a preflight request to
	// indicate which HTTP headers can be used when making the actual request.
	//
	// Optional. Default value []string{}.
	//
	// See also: https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Allow-Headers
	AllowHeaders []string `json:"allow_headers,omitempty" yaml:"allow_headers,omitempty"`

	// AllowCredentials determines the value of the
	// Access-Control-Allow-Credentials response header.  This header indicates
	// whether or not the response to the request can be exposed when the
	// credentials mode (Request.credentials) is true. When used as part of a
	// response to a preflight request, this indicates whether or not the actual
	// request can be made using credentials.  See also
	// [MDN: Access-Control-Allow-Credentials].
	//
	// Optional. Default value false, in which case the header is not set.
	//
	// Security: avoid using `AllowCredentials = true` with `AllowOrigins = *`.
	// See "Exploiting CORS misconfigurations for Bitcoins and bounties",
	// https://blog.portswigger.net/2016/10/exploiting-cors-misconfigurations-for.html
	//
	// See also: https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Allow-Credentials
	AllowCredentials bool `json:"allow_credentials,omitempty" yaml:"allow_credentials,omitempty"`

	// UnsafeWildcardOriginWithAllowCredentials UNSAFE/INSECURE: allows wildcard '*' origin to be used with AllowCredentials
	// flag. In that case we consider any origin allowed and send it back to the client with `Access-Control-Allow-Origin` header.
	//
	// This is INSECURE and potentially leads to [cross-origin](https://portswigger.net/research/exploiting-cors-misconfigurations-for-bitcoins-and-bounties)
	// attacks. See: https://github.com/labstack/echo/issues/2400 for discussion on the subject.
	//
	// Optional. Default value is false.
	UnsafeWildcardOriginWithAllowCredentials bool `json:"unsafe_wildcard_origin_with_allow_credentials,omitempty" yaml:"unsafe_wildcard_origin_with_allow_credentials,omitempty"`

	// ExposeHeaders determines the value of Access-Control-Expose-Headers, which
	// defines a list of headers that clients are allowed to access.
	//
	// Optional. Default value []string{}, in which case the header is not set.
	//
	// See also: https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Expose-Header
	ExposeHeaders []string `json:"expose_headers,omitempty" yaml:"expose_headers,omitempty"`

	// MaxAge determines the value of the Access-Control-Max-Age response header.
	// This header indicates how long (in seconds) the results of a preflight
	// request can be cached.
	// The header is set only if MaxAge != 0, negative value sends "0" which instructs browsers not to cache that response.
	//
	// Optional. Default value 0 - meaning header is not sent.
	//
	// See also: https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Max-Age
	MaxAge int `json:"max_age,omitempty" yaml:"max_age,omitempty"`
}

type CSRFConfig struct {
	// Skipper defines a function to skip middleware.
	Skipper middleware.Skipper `json:"-" yaml:"-"`

	// ErrorHandler defines a function which is executed for returning custom errors.
	ErrorHandler middleware.CSRFErrorHandler `json:"-" yaml:"-"`

	// TokenLength is the length of the generated token.
	TokenLength uint8 `json:"token_length,omitempty" yaml:"token_length,omitempty"`
	// Optional. Default value 32.

	// TokenLookup is a string in the form of "<source>:<name>" or "<source>:<name>,<source>:<name>" that is used
	// to extract token from the request.
	// Optional. Default value "header:X-CSRF-Token".
	// Possible values:
	// - "header:<name>" or "header:<name>:<cut-prefix>"
	// - "query:<name>"
	// - "form:<name>"
	// Multiple sources example:
	// - "header:X-CSRF-Token,query:csrf"
	TokenLookup string `json:"token_lookup,omitempty" yaml:"token_lookup,omitempty"`

	// Context key to store generated CSRF token into context.
	// Optional. Default value "csrf".
	ContextKey string `json:"context_key,omitempty" yaml:"context_key,omitempty"`

	Cookie struct {
		Name     string        `json:"name" yaml:"name"`
		Domain   string        `json:"domain" yaml:"domain"`
		Path     string        `json:"path" yaml:"path"`
		MaxAge   time.Duration `json:"max_age,omitempty" yaml:"max_age,omitempty"`
		Secure   bool          `json:"secure" yaml:"secure"`
		HTTPOnly bool          `json:"http_only" yaml:"http_only"`
		SameSite SameSiteType  `json:"same_site" yaml:"same_site"`
	} `json:"cookie" yaml:"cookie"`
}

type JWTConfig struct {
	Secret               string        `json:"secret,omitempty" yaml:"secret,omitempty"`
	AccessTokenDuration  time.Duration `json:"access_token_duration,omitempty" yaml:"access_token_duration,omitempty"`
	RefreshTokenDuration time.Duration `json:"refresh_token_duration,omitempty" yaml:"refresh_token_duration,omitempty"`
}

func (cfg *JWTConfig) InitDefaults() {
	if cfg.Secret == "" {
		cfg.Secret = cms.RandomString(50)
	}
	if cfg.AccessTokenDuration == 0 {
		cfg.AccessTokenDuration = 5 * time.Minute
	}
	if cfg.RefreshTokenDuration == 0 {
		cfg.RefreshTokenDuration = 60 * time.Minute
	}
}
