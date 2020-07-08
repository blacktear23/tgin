package tgin

import (
	"net/http"
)

var (
	_ http.ResponseWriter = (*ResponseWriterWrapper)(nil)
)

type H map[string]interface{}

type ResponseWriterWrapper struct {
	writer http.ResponseWriter
	code   int
	ctx    *Context
}

func (w *ResponseWriterWrapper) Header() http.Header {
	return w.writer.Header()
}

func (w *ResponseWriterWrapper) Write(b []byte) (int, error) {
	return w.writer.Write(b)
}

func (w *ResponseWriterWrapper) WriteHeader(code int) {
	w.code = code
	w.writer.WriteHeader(code)
}
