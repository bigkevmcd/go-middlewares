package logging

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-logr/logr"
)

// RequestLogger sets up a post-request handler to log out the request.
func RequestLogger(l logr.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			wrapped := WrapWriter(w)
			start := time.Now()
			defer func() {
				rl := WithRequestValues(l, r)
				rl.WithValues("duration", fmt.Sprintf("%dms", time.Since(start).Milliseconds())).Info(http.StatusText(wrapped.Status()))
			}()
			next.ServeHTTP(wrapped, r)
		}
		return http.HandlerFunc(fn)
	}
}
