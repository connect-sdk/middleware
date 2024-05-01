package middleware

import (
	"net/http"
	"net/url"
)

// GetRequestURL returns the request URL
func GetRequestURL(r *http.Request) *url.URL {
	endpoint := &url.URL{}

	if r.TLS == nil {
		endpoint.Scheme = "http"
	} else {
		endpoint.Scheme = "https"
	}

	if r.Host == "" {
		endpoint.Host = r.Header.Get("Host")
	} else {
		endpoint.Host = r.Host
	}

	if value := r.Header.Get("X-Forwarded-Host"); value != "" {
		endpoint.Host = value
	}

	if value := r.Header.Get("X-Forwarded-Proto"); value != "" {
		endpoint.Scheme = value
	}

	if value := r.URL; value != nil {
		endpoint.Path = value.Path
	}

	return endpoint
}
