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
		return http.HandlerFunc(wrapResponseWithLogger(l, next))
	}
}

// TODO: refactor with a predicate?
//
// allow passing in a predicate function func(*http.Request) bool
// impact on performance?
func wrapResponseWithLogger(l logr.Logger, next http.Handler) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		wrapped := WrapWriter(w)
		start := time.Now()
		defer func() {
			rl := WithRequestValues(l, r)
			rl.WithValues("duration", fmt.Sprintf("%dms", time.Since(start).Milliseconds())).Info(http.StatusText(wrapped.Status()))
		}()
		next.ServeHTTP(wrapped, r)
	}
}
