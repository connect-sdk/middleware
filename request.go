package middleware

import (
	"net/http"
	"net/url"
)

// GetRequestURL returns the request URL
func GetRequestURL(r *http.Request) *url.URL {
	addr := &url.URL{}

	if r.TLS == nil {
		addr.Scheme = "http"
	} else {
		addr.Scheme = "https"
	}

	if r.Host == "" {
		addr.Host = r.Header.Get("Host")
	} else {
		addr.Host = r.Host
	}

	if value := r.Header.Get("X-Forwarded-Host"); value != "" {
		addr.Host = value
	}

	if value := r.Header.Get("X-Forwarded-Proto"); value != "" {
		addr.Scheme = value
	}

	return addr
}
