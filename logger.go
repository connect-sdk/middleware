package middleware

import (
	"net/http"
	"time"

	"github.com/ralch/slogr"
)

// WithLogger set up the http logger.
func WithLogger() func(http.Handler) http.Handler {
	innerFn := func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			start := time.Now()

			ww := &ResponseWriter{ResponseWriter: w}
			// prepare the logger
			logger := slogr.FromContext(ctx)
			logger = logger.With(slogr.Request(r))

			// prepare the context
			ctx = slogr.WithContext(ctx, logger)
			// prepare the request
			r = r.WithContext(ctx)

			// execute the handler
			next.ServeHTTP(ww, r)

			duration := time.Since(start)
			// log the request end
			logger = logger.With(slogr.ResponseWriter(ww, slogr.WithLatency(duration)))

			status := ww.GetStatusCode()
			switch {
			case status < 400:
				logger.InfoContext(ctx, "")
			case status < 500:
				logger.WarnContext(ctx, "")
			default:
				logger.ErrorContext(ctx, "")
			}
		}

		return http.HandlerFunc(fn)
	}

	return innerFn
}

var (
	_ http.Flusher        = &ResponseWriter{}
	_ http.ResponseWriter = &ResponseWriter{}
)

// ResponseWriter repersents a response writer.
type ResponseWriter struct {
	ResponseWriter http.ResponseWriter
	StatusCode     int32
	ContentLength  int64
}

// Flush implements http.Flusher.
func (r *ResponseWriter) Flush() {
	if flusher, ok := r.ResponseWriter.(http.Flusher); ok {
		flusher.Flush()
	}
}

// Header implements http.ResponseWriter
func (r *ResponseWriter) Header() http.Header {
	return r.ResponseWriter.Header()
}

// Write implements http.ResponseWriter
func (r *ResponseWriter) Write(data []byte) (int, error) {
	n, err := r.ResponseWriter.Write(data)
	r.ContentLength = r.ContentLength + int64(n)
	return n, err
}

// WriteHeader implements http.ResponseWriter
func (r *ResponseWriter) WriteHeader(code int) {
	r.StatusCode = int32(code)
	r.ResponseWriter.WriteHeader(code)
}

// GetStatusCode returns the StatusCode.
func (r *ResponseWriter) GetStatusCode() int32 {
	return r.StatusCode
}

// GetContentLength returns the ContentLength.
func (r *ResponseWriter) GetContentLength() int64 {
	return r.ContentLength
}
