package tgin

import (
	"net/http"
)

var (
	_ http.ResponseWriter = (*ResponseWriterWrapper)(nil)
	_ http.Hijacker       = (*ResponseWriterWrapper)(nil)
)

type H map[string]interface{}

type ResponseWriterWrapper struct {
	http.ResponseWriter
	http.Hijacker
	code int
	ctx  *Context
}

func (w *ResponseWriterWrapper) WriteHeader(code int) {
	w.code = code
	w.ResponseWriter.WriteHeader(code)
}
