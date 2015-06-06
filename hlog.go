package hlog

import (
	"bufio"
	"log"
	"net"
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

func (w *statusLoggingResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	hijacker := w.ResponseWriter.(http.Hijacker)
	return hijacker.Hijack()
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
