package middleware

import (
	"compress/gzip"
	"net/http"
	"strings"
)

// GzipHandler compresses the responses to send files.
type GzipHandler struct {
}

// WrappedReponseWriter encapsules a compressed response object.
type WrappedReponseWriter struct {
	rw http.ResponseWriter
	gw *gzip.Writer
}

// GzipMiddleware intercepts the response to compress the contents.
func (g *GzipHandler) GzipMiddleware(next http.Handler) http.Handler {

	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			// create a gziped response
			wrw := NewWrappedResponseWriter(rw)
			wrw.Header().Set("Content-Encoding", "gzip")

			next.ServeHTTP(wrw, r)
			defer wrw.Flush()

			return
		}

		// handle normal
		next.ServeHTTP(rw, r)
	})
}

// NewWrappedResponseWriter creates a new WrappedReponseWriter
func NewWrappedResponseWriter(rw http.ResponseWriter) *WrappedReponseWriter {
	gw := gzip.NewWriter(rw)

	return &WrappedReponseWriter{rw: rw, gw: gw}
}

// Header sets the response header.
func (wr *WrappedReponseWriter) Header() http.Header {
	return wr.rw.Header()
}

// Write writes the contents to the response.
func (wr *WrappedReponseWriter) Write(d []byte) (int, error) {
	return wr.gw.Write(d)
}

// WriteHeader writes the status code to the response.
func (wr *WrappedReponseWriter) WriteHeader(statuscode int) {
	wr.rw.WriteHeader(statuscode)
}

// Flush flushes the response.
func (wr *WrappedReponseWriter) Flush() {
	wr.gw.Flush()
	wr.gw.Close()
}
