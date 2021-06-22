package logging

import (
	"net/http"
)

// ResponseWriterWrapper is a proxy around an http.ResponseWriter that allows you to hook
// into various parts of the response process.
type ResponseWriterWrapper interface {
	http.ResponseWriter
	// Status returns the HTTP status of the request, or 0 if one has not
	// yet been sent.
	Status() int
	// BytesWritten returns the total number of bytes sent to the client.
	BytesWritten() int
	// Unwrap returns the original proxied target.
	Unwrap() http.ResponseWriter
}

// WrapWriter wraps an http.ResponseWriter in a proxy that implements the
// ResponseWriterWrapper interface.
func WrapWriter(w http.ResponseWriter) ResponseWriterWrapper {
	return &writerWrapper{ResponseWriter: w}
}

type writerWrapper struct {
	http.ResponseWriter
	code        int
	bytes       int
	wroteHeader bool
}

// WriteHeader replaces the default ResponseWriter implementation.
func (w *writerWrapper) WriteHeader(code int) {
	if !w.wroteHeader {
		w.code = code
		w.wroteHeader = true
		w.ResponseWriter.WriteHeader(code)
	}
}

// Write replaces the default ResponseWriter implementation.
func (w *writerWrapper) Write(buf []byte) (int, error) {
	if !w.wroteHeader {
		w.WriteHeader(http.StatusOK)
	}
	n, err := w.ResponseWriter.Write(buf)
	w.bytes += n
	return n, err
}

// Status is an implementation of the ResponseWriterWrapper interface.
func (w *writerWrapper) Status() int {
	return w.code
}

// BytesWritten is an implementation of the ResponseWriterWrapper interface.
func (w *writerWrapper) BytesWritten() int {
	return w.bytes
}

// Unwrap is an implementation of the ResponseWriterWrapper interface.
func (w *writerWrapper) Unwrap() http.ResponseWriter {
	return w.ResponseWriter
}
