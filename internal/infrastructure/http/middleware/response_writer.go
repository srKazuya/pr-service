package middleware

import "net/http"

type responseWriter struct {
	http.ResponseWriter
	status      int
	bytesWritten int
}

func NewWrapResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{ResponseWriter: w, status: 200}
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	n, err := rw.ResponseWriter.Write(b)
	rw.bytesWritten += n
	return n, err
}

func (rw *responseWriter) Status() int {
	return rw.status
}

func (rw *responseWriter) BytesWritten() int {
	return rw.bytesWritten
}
