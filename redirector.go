package middleware

import (
	"bytes"
	"net/http"
)

// WithRedirector represents an method redirctor
func WithRedirector() func(http.Handler) http.Handler {
	innerFn := func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			writer := &ResponseWriter{
				ResponseWriter: w,
				ResponseBody:   &bytes.Buffer{},
			}

			next.ServeHTTP(writer, r)

			if location := writer.Header().Get("X-Location"); location != "" {
				http.Redirect(w, r, location, http.StatusFound)
			} else {
				writer.Flush()
			}
		}

		return http.HandlerFunc(fn)
	}

	return innerFn
}
