package logging

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWrappedWriter(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusAccepted)
		io.WriteString(w, "This is a test body of exactly 39 bytes")
	}
	req := httptest.NewRequest("GET", "http://example.com/foo", nil)
	rec := httptest.NewRecorder()
	w := WrapWriter(rec)

	handler(w, req)

	if s := w.Status(); s != http.StatusAccepted {
		t.Fatalf("Status() got %d, want %d", s, http.StatusAccepted)
	}

	if b := w.BytesWritten(); b != 39 {
		t.Fatalf("BytesWritten() got %d, want %d", b, 39)
	}

	if tw := w.Unwrap(); tw != rec {
		t.Fatalf("Unwrap() got %v, want %v", tw, rec)
	}
}

func TestWrappedWriter_defaults_status_to_ok(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "This is a body of exactly 34 bytes")
	}
	req := httptest.NewRequest("GET", "http://example.com/foo", nil)
	rec := httptest.NewRecorder()
	w := WrapWriter(rec)

	handler(w, req)

	if s := w.Status(); s != http.StatusOK {
		t.Fatalf("Status() got %d, want %d", s, http.StatusOK)
	}

	if b := w.BytesWritten(); b != 34 {
		t.Fatalf("BytesWritten() got %d, want %d", b, 34)
	}
}

func TestWrappedWriter_accumulates_writes_written(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "This is 16 bytes")
		io.WriteString(w, "This is 16 bytes")
	}
	req := httptest.NewRequest("GET", "http://example.com/foo", nil)
	rec := httptest.NewRecorder()
	w := WrapWriter(rec)

	handler(w, req)

	if b := w.BytesWritten(); b != 32 {
		t.Fatalf("BytesWritten() got %d, want %d", b, 32)
	}
}

func TestWrappedWriter_only_writes_status_once(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusAccepted)
		w.WriteHeader(http.StatusCreated)
		io.WriteString(w, "This is 16 bytes")
	}
	req := httptest.NewRequest("GET", "http://example.com/foo", nil)
	rec := httptest.NewRecorder()
	w := WrapWriter(rec)

	handler(w, req)

	if s := w.Status(); s != http.StatusAccepted {
		t.Fatalf("Status() got %d, want %d", s, http.StatusAccepted)
	}
}
