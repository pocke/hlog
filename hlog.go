package hlog

import (
	"log"
	"net/http"
)

type statusLoggingResponseWriter struct {
	http.ResponseWriter
	status int
}

func (w *statusLoggingResponseWriter) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

// check interface
var _ http.ResponseWriter = &statusLoggingResponseWriter{}

func Wrap(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		wLog := &statusLoggingResponseWriter{w, 200}

		log.Printf(`Started %s "%s" for %s`, r.Method, r.RequestURI, r.RemoteAddr)
		f(wLog, r)
		log.Printf("response status: %d", wLog.status)
	}
}
