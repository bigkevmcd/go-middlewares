package logging

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-logr/logr"
	"github.com/go-logr/logr/funcr"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestOptionalRequestLogger_no_param(t *testing.T) {
	var logged []map[string]string
	logger := funcr.New(func(_, args string) {
		logged = append(logged, parseArgs(args))
	}, funcr.Options{})

	handler := OptionalRequestLogger(logger)(http.HandlerFunc(testOKHandler(t)))
	runRequestLoggerTestWithPath(t, logger, handler, "")

	var want []map[string]string
	if diff := cmp.Diff(want, logged, cmpopts.IgnoreMapEntries(func(k string, v interface{}) bool {
		return k == "host" || k == "remote_addr"
	})); diff != "" {
		t.Fatalf("failed to capture logs:\n%s", diff)
	}
}

func TestOptionalRequestLogger_with_param(t *testing.T) {
	var logged []map[string]string
	logger := funcr.New(func(_, args string) {
		logged = append(logged, parseArgs(args))
	}, funcr.Options{})

	handler := OptionalRequestLogger(logger)(http.HandlerFunc(testOKHandler(t)))
	runRequestLoggerTestWithPath(t, logger, handler, "/test?lg")

	want := []map[string]string{
		{
			"level":      "0",
			"method":     "GET",
			"msg":        "OK",
			"path":       "/test?lg",
			"user_agent": "Go-http-client/1.1",
			"duration":   "0ms",
			"protocol":   "",
		},
	}
	if diff := cmp.Diff(want, logged, cmpopts.IgnoreMapEntries(func(k string, v interface{}) bool {
		return k == "host" || k == "remote_addr"
	})); diff != "" {
		t.Fatalf("failed to capture logs:\n%s", diff)
	}
}

func runRequestLoggerTestWithPath(t *testing.T, logger logr.Logger, h http.Handler, urlPath string) {
	t.Helper()
	s := httptest.NewServer(h)
	defer s.Close()

	resp, err := http.Get(s.URL + urlPath)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("got StatusCode %v, want %v", resp.StatusCode, http.StatusOK)
	}

	if data, _ := ioutil.ReadAll(resp.Body); string(data) != "ok" {
		t.Fatalf("Body = %v, want %v", string(data), "ok")
	}
}
