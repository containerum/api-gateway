package middleware

import (
	"net/http"
)

// WrapResponseWriter is a proxy around an http.ResponseWriter that allows you to hook
// into various parts of the response process.
type WrapResponseWriter interface {
	http.ResponseWriter
	// Status returns the HTTP status of the request, or 0 if one has not
	// yet been sent.
	Status() int
	// BytesWritten returns the total number of bytes sent to the client.
	BytesWritten() int
}

// NewWrapResponseWriter wraps an http.ResponseWriter, returning a proxy that allows you to
// hook into various parts of the response process.
func NewWrapResponseWriter(w http.ResponseWriter) WrapResponseWriter {
	bw := basicWriter{ResponseWriter: w}
	return &bw
}

// basicWriter wraps a http.ResponseWriter that implements the minimal
// http.ResponseWriter interface.
type basicWriter struct {
	http.ResponseWriter
	wroteHeader bool
	code        int
	bytes       int
	headers     []string
}

func (b *basicWriter) WriteHeader(code int) {
	if !b.wroteHeader {
		b.code = code
		b.wroteHeader = true
		b.ResponseWriter.WriteHeader(code)
	}
}

func (b *basicWriter) Write(buf []byte) (int, error) {
	b.WriteHeader(http.StatusOK)
	n, err := b.ResponseWriter.Write(buf)
	b.bytes += n
	return n, err
}

func (b *basicWriter) Status() int {
	return b.code
}

func (b *basicWriter) BytesWritten() int {
	return b.bytes
}

// func (b *basicWriter) Header() http.Header {
// 	http.He
// }
