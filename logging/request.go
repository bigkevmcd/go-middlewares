package logging

import (
	"net/http"

	"github.com/go-logr/logr"
)

// WithRequestValues adds various items from the request to the values associated with
// the logger for use in logging HTTP requests.
func WithRequestValues(l logr.Logger, r *http.Request) logr.Logger {
	return l.WithValues(
		"method", r.Method,
		"host", r.Host,
		"path", r.URL.RequestURI(),
		"remote_addr", r.RemoteAddr,
		"user_agent", r.UserAgent(),
		"protocol", r.URL.Scheme,
	)
}
