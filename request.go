package middleware

import (
	"net/http"
	"net/url"
)

// GetRequestURL returns the request URL
func GetRequestURL(r *http.Request) *url.URL {
	uri := &url.URL{
		Host: r.Header.Get("Host"),
	}

	if r.TLS == nil {
		uri.Scheme = "http"
	} else {
		uri.Scheme = "https"
	}

	if value := r.Header.Get("X-Forwarded-Host"); value != "" {
		uri.Host = value
	}

	if value := r.Header.Get("X-Forwarded-Proto"); value != "" {
		uri.Scheme = value
	}

	return uri
}
