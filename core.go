package middleware

import "net/http"

// RoundTripperOption is an option for configuring a RoundTripper.
type RoundTripperOption interface {
	// WrapRoundTripper wraps the given RoundTripper with some custom logic.
	WrapRoundTripper(http.RoundTripper) http.RoundTripper
}

var _ RoundTripperOption = RoundTripperOptionFunc(nil)

// RoundTripperOptionFunc is a function type that implements RoundTripperOption.
type RoundTripperOptionFunc func(http.RoundTripper) http.RoundTripper

// WrapRoundTripper implements RoundTripperOption.
func (fn RoundTripperOptionFunc) WrapRoundTripper(r http.RoundTripper) http.RoundTripper {
	return fn(r)
}
