package tgin

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type Context struct {
	Method  string
	Request *http.Request
	Writer  http.ResponseWriter
	aborted bool
	values  map[string]interface{}
}

func newContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		Method:  r.Method,
		Request: r,
		Writer:  w,
		aborted: false,
		values:  make(map[string]interface{}),
	}
}

func (c *Context) json(code int, val H, indented bool) {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	if indented {
		enc.SetIndent("", "    ")
	}
	err := enc.Encode(val)
	if err != nil {
		c.Text(500, fmt.Sprintf("Server Error!\n%v", err))
		return
	}
	w := c.Writer
	hdr := w.Header()
	hdr.Set("Content-Type", "application/json; charset=utf-8")
	hdr.Set("Content-Length", fmt.Sprintf("%d", buf.Len()))
	w.WriteHeader(code)
	buf.WriteTo(w)
}

func (c *Context) BindJSON(obj interface{}) error {
	dec := json.NewDecoder(c.Request.Body)
	return dec.Decode(obj)
}

func (c *Context) JSON(code int, val H) {
	c.json(code, val, false)
}

func (c *Context) IndentedJSON(code int, val H) {
	c.json(code, val, true)
}

func (c *Context) Text(code int, message string) {
	w := c.Writer
	hdr := w.Header()
	hdr.Set("Content-Type", "text/plain; charset=utf-8")
	hdr.Set("Content-Length", fmt.Sprintf("%d", len(message)))
	w.WriteHeader(code)
	w.Write([]byte(message))
}

func (c *Context) GetQuery(key string) (string, bool) {
	ret := c.Request.URL.Query().Get(key)
	if ret == "" {
		return ret, false
	}
	return ret, true
}

func (c *Context) GetHeader(key string) string {
	return c.Request.Header.Get(key)
}

func (c *Context) File(filePath string) {
	http.ServeFile(c.Writer, c.Request, filePath)
}

func (c *Context) Abort() {
	c.aborted = true
}

func (c *Context) AbortWithStatus(code int) {
	w := c.Writer
	w.WriteHeader(code)
	c.Abort()
}

func (c *Context) Set(key string, obj interface{}) {
	c.values[key] = obj
}

func (c *Context) Get(key string) (interface{}, bool) {
	obj, have := c.values[key]
	return obj, have
}
