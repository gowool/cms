package http

import "net/http"

const (
	headerXForwardedProto    = "X-Forwarded-Proto"
	headerXForwardedProtocol = "X-Forwarded-Protocol"
	headerXForwardedSsl      = "X-Forwarded-Ssl"
	headerXUrlScheme         = "X-Url-Scheme"
)

func Scheme(r *http.Request) string {
	if r.TLS != nil {
		return "https"
	}
	if scheme := r.Header.Get(headerXForwardedProto); scheme != "" {
		return scheme
	}
	if scheme := r.Header.Get(headerXForwardedProtocol); scheme != "" {
		return scheme
	}
	if ssl := r.Header.Get(headerXForwardedSsl); ssl == "on" {
		return "https"
	}
	if scheme := r.Header.Get(headerXUrlScheme); scheme != "" {
		return scheme
	}
	return "http"
}
