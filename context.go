package tgin

import (
	"bytes"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
)

const formMaxMemory = 32 << 20 // 32 MB

type Context struct {
	Method      string
	Request     *http.Request
	Writer      http.ResponseWriter
	aborted     bool
	served      bool
	values      map[string]interface{}
	middlewares []RouteHandler
	index       int
	mux         *http.ServeMux
}

func newContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		Method:      r.Method,
		Request:     r,
		Writer:      w,
		aborted:     false,
		served:      false,
		index:       0,
		middlewares: nil,
		values:      make(map[string]interface{}),
	}
}

func (c *Context) json(code int, val interface{}, indented bool) {
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

func (c *Context) JSON(code int, val interface{}) {
	c.json(code, val, false)
}

func (c *Context) IndentedJSON(code int, val interface{}) {
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

func (c *Context) String(code int, format string, values ...interface{}) {
	body := fmt.Sprintf(format, values...)
	c.Text(code, body)
}

func (c *Context) GetQuery(key string) (string, bool) {
	ret := c.Request.URL.Query().Get(key)
	if ret == "" {
		return ret, false
	}
	return ret, true
}

func (c *Context) GetQueryArray(key string) ([]string, bool) {
	val := c.Request.URL.Query()[key]
	if len(val) == 0 {
		return []string{}, false
	}
	return val, true
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

func (c *Context) Header(key, value string) {
	if value == "" {
		c.Writer.Header().Del(key)
		return
	}
	c.Writer.Header().Set(key, value)
}

func (c *Context) Redirect(code int, location string) {
	http.Redirect(c.Writer, c.Request, location, code)
}

func (c *Context) PostForm(key string) string {
	ret := c.PostFormArray(key)
	if len(ret) == 0 {
		return ""
	}
	return ret[0]
}

func (c *Context) PostFormArray(key string) []string {
	if c.Request.Form == nil {
		c.Request.ParseMultipartForm(formMaxMemory)
	}
	vals := c.Request.Form[key]
	if len(vals) == 0 {
		return []string{}
	}
	return vals
}

func (c *Context) FormFile(key string) (*multipart.FileHeader, error) {
	_, header, err := c.Request.FormFile(key)
	return header, err
}

func (c *Context) Next() {
	numMw := len(c.middlewares)
	if c.index < numMw {
		for c.index < numMw {
			m := c.middlewares[c.index]
			c.index++
			m(c)
			if c.aborted {
				break
			}
		}
	}
	// Middleware execute all We should serve it
	if !c.served && !c.aborted {
		c.served = true
		c.mux.ServeHTTP(c.Writer, c.Request)
	}
}
