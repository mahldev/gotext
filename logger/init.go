package logger

import (
	"log"
	"net/http"
)

type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}

func Logger(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		method := r.Method
		url := r.URL.String()
		remoteAddr := r.RemoteAddr

		log.Printf("[INFO] Incoming request: method=%s, URL=%s, remoteAddr=%s\n", method, url, remoteAddr)

		rw := &responseWriter{w, http.StatusOK}

		next.ServeHTTP(rw, r)
		statusCode := rw.status

		log.Printf("[INFO] Request completed: status=%d\n", statusCode)
	}

	return http.HandlerFunc(fn)
}
