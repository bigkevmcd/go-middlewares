package logging

import (
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/go-logr/logr/funcr"
	"github.com/google/go-cmp/cmp"
)

func parseArgs(s string) map[string]string {
	dequote := func(s string) string {
		return strings.Trim(s, "\"")
	}
	splitKV := func(kv string) (string, string) {
		x := strings.SplitN(kv, "=", 2)
		return dequote(x[0]), dequote(x[1])
	}

	args := make(map[string]string)
	items := strings.Split(s, " ")
	for _, v := range items {
		key, value := splitKV(v)
		args[key] = value
	}
	return args
}

func TestWithRequestValues(t *testing.T) {
	var logged []map[string]string
	logger := funcr.New(func(_, args string) {
		logged = append(logged, parseArgs(args))
	}, funcr.Options{})

	r := mustNewRequest(t, http.MethodGet, "https://example.com/testing/", nil, func(r *http.Request) {
		r.RemoteAddr = "192.168.0.1:5700"
		r.Header.Set("User-Agent", "test-agent/0.1")
	})

	l := WithRequestValues(logger, r)

	l.Info("testing")

	want := []map[string]string{
		{
			"host":        "example.com",
			"level":       "0",
			"method":      "GET",
			"msg":         "testing",
			"path":        "/testing/",
			"protocol":    "https",
			"remote_addr": "192.168.0.1:5700",
			"user_agent":  "test-agent/0.1",
		},
	}
	if diff := cmp.Diff(want, logged); diff != "" {
		t.Fatalf("failed to capture logs:\n%s", diff)
	}
}

func mustNewRequest(t *testing.T, method, url string, body io.Reader, opts ...func(*http.Request)) *http.Request {
	t.Helper()
	r, err := http.NewRequest(method, url, body)
	if err != nil {
		t.Fatal(err)
	}
	for _, o := range opts {
		o(r)
	}
	return r
}
