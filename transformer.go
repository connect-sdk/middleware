package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"mime"
	"net/http"
	"strconv"
)

// WithTransformer transforms the input.
func WithTransformer() func(http.Handler) http.Handler {
	decoder := NewFormDecoder()

	innerFn := func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			mime, _, _ := mime.ParseMediaType(r.Header.Get("Content-Type"))

			if mime == "application/x-www-form-urlencoded" {
				if err := r.ParseForm(); err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}

				if err := r.ParseForm(); err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}

				// find the message by name
				factory, err := URLMessageType(r.URL)
				if err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}

				message := factory.New().Interface()
				// decode the message
				if err = decoder.Decode(message, r.Form); err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}

				body := &bytes.Buffer{}
				// encode the boyd
				if err := json.NewEncoder(body).Encode(message); err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}

				r.Body = io.NopCloser(body)
				r.ContentLength = int64(body.Len())

				r.Header.Set("Content-Type", "application/json")
				r.Header.Set("Content-Length", strconv.Itoa(body.Len()))
			}

			// execute the handler
			next.ServeHTTP(w, r)
		}

		return http.HandlerFunc(fn)
	}

	return innerFn
}
