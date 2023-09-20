package middleware

import (
	"net/http"

	"github.com/connect-sdk/telemetry"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

// WithTracer represents an operation telemetry
func WithTracer() func(http.Handler) http.Handler {
	propagator := telemetry.NewTracePropagator()

	innerFn := func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			options := []otelhttp.Option{
				// setup the propagators
				otelhttp.WithPropagators(
					propagator,
				),
				// setup the events
				otelhttp.WithMessageEvents(
					otelhttp.ReadEvents,
					otelhttp.WriteEvents,
				),
			}
			// setup the handler
			handler := otelhttp.NewHandler(next, r.URL.Path, options...)
			handler.ServeHTTP(w, r)
		}

		return http.HandlerFunc(fn)
	}

	return innerFn
}
