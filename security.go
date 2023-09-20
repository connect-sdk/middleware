package middleware

import (
	"context"
	"net/http"

	"github.com/go-chi/cors"
	"google.golang.org/api/option"
	htransport "google.golang.org/api/transport/http"
)

// CorsConfig represents the CORS options.
type CorsConfig struct {
	// AllowedOrigins is a list of origins a cross-domain request can be executed from.
	// If the special "*" value is present in the list, all origins will be allowed.
	// An origin may contain a wildcard (*) to replace 0 or more characters
	// (i.e.: http://*.domain.com). Usage of wildcards implies a small performance penalty.
	// Only one wildcard can be used per origin.
	// Default value is ["*"]
	AllowedOrigins []string
	// AllowedMethods is a list of methods the client is allowed to use with cross-domain requests. Default value is simple methods (HEAD, GET and POST).
	AllowedMethods []string
	// AllowedHeaders is list of non simple headers the client is allowed to use with
	// cross-domain requests.
	// If the special "*" value is present in the list, all headers will be allowed.
	// Default value is [] but "Origin" is always appended to the list.
	AllowedHeaders []string
	// ExposedHeaders indicates which headers are safe to expose to the API of a CORS API specification
	ExposedHeaders []string
}

// WithCors creates a new Cors handler with passed options.
func WithCors(config ...CorsConfig) func(http.Handler) http.Handler {
	innterFn := func(next http.Handler) http.Handler {
		options := cors.Options{
			AllowedOrigins: []string{"*"},
			AllowedMethods: []string{
				http.MethodGet,
				http.MethodPost,
			},
			AllowedHeaders: []string{
				"Authorization",            // Used by HTTP
				"Accept-Encoding",          // Used by HTTP
				"Content-Encoding",         // Used by HTTP
				"Content-Type",             // Used by HTTP
				"Connect-Protocol-Version", // Used by Connect
				"Connect-Timeout-Ms",       // Used by Connect
				"Connect-Accept-Encoding",  // Unused in web browsers, but added for future-proofing
				"Connect-Content-Encoding", // Unused in web browsers, but added for future-proofing
				"Grpc-Timeout",             // Used for gRPC-web
				"X-Grpc-Web",               // Used for gRPC-web
				"X-User-Agent",             // Used for gRPC-web
				"X-Request-ID",             // Used by HTTP
				"X-Requested-With",         // Used by HTTP
				"X-HTTP-Method-Override",   // Used by HTTP
			},
			ExposedHeaders: []string{
				"Content-Encoding",         // Unused in web browsers, but added for future-proofing
				"Connect-Content-Encoding", // Unused in web browsers, but added for future-proofing
				"Grpc-Status",              // Required for gRPC-web
				"Grpc-Message",             // Required for gRPC-web
				"Location",                 // Used by HTTP
			},
			AllowCredentials: true,
			MaxAge:           86400,
		}

		for _, cfg := range config {
			options.AllowedOrigins = append(options.AllowedOrigins, cfg.AllowedOrigins...)
			options.AllowedMethods = append(options.AllowedMethods, cfg.AllowedMethods...)
			options.AllowedHeaders = append(options.AllowedHeaders, cfg.AllowedHeaders...)
			options.ExposedHeaders = append(options.ExposedHeaders, cfg.ExposedHeaders...)
		}

		// taken from https://connectrpc.com/docs/go/deployment#cors
		fn := cors.Handler(options)

		return fn(next)
	}

	return innterFn
}

// TransportWithAuthorization configures the transport to use authorize the client with given audience.
func TransportWithAuthorization(audience string, options ...option.ClientOption) RoundTripperOption {
	innerFn := func(transport http.RoundTripper) http.RoundTripper {
		if transport == nil {
			transport = http.DefaultTransport
		}

		transport, err := htransport.NewTransport(context.Background(), transport, options...)
		if err != nil {
			panic(err)
		}

		return transport
	}

	return RoundTripperOptionFunc(innerFn)
}
