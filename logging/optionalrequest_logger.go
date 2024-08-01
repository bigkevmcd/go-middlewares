package logging

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-logr/logr"
)

// OptionalRequestLogger sets up a post-request handler to log out the request.
//
// The Request is only logged out if the URL includes a query param ?lg
// TODO: Configurable query param?
// TODO: Also allow a specific header?
func OptionalRequestLogger(l logr.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Query().Has("lg") {
				wrapped := WrapWriter(w)
				start := time.Now()
				defer func() {
					rl := WithRequestValues(l, r)
					rl.WithValues("duration", fmt.Sprintf("%dms", time.Since(start).Milliseconds())).Info(http.StatusText(wrapped.Status()))
				}()
				next.ServeHTTP(wrapped, r)
			} else {
				next.ServeHTTP(w, r)
			}
		}

		return http.HandlerFunc(fn)
	}
}
