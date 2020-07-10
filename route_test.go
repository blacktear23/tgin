package tgin

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func processRequest(r *RouteGroup, method, path string, body ...string) *http.Response {
	var buf *bytes.Buffer
	if len(body) == 0 {
		buf = bytes.NewBuffer(nil)
	} else {
		buf = bytes.NewBufferString(strings.Join(body, " "))
	}
	req := httptest.NewRequest(method, path, buf)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w.Result()
}

func TestGetHandler(t *testing.T) {
	r := NewRouteGroup()
	r.Get("/hello", func(c *Context) {
		c.String(200, "Hello World")
	})
	resp := processRequest(r, "GET", "/hello")
	assertEqual(t, 200, resp.StatusCode, "Status Code Error")
	body := ReadBodyString(resp)
	assertEqual(t, "Hello World", body, "Body not correct")

	resp = processRequest(r, "HEAD", "/hello")
	assertEqual(t, 404, resp.StatusCode, "Status Code Error")
}

func TestUseMiddlewareTransferObject(t *testing.T) {
	r := NewRouteGroup()
	r.Use(func(c *Context) {
		c.Set("test", "Test")
	})
	r.Get("/hello", func(c *Context) {
		val, _ := c.Get("test")
		assertEqual(t, "Test", val.(string))
		c.String(200, "OK")
	})
	resp := processRequest(r, "GET", "/hello")
	assertEqual(t, 200, resp.StatusCode, "Status Code Error")
}

func TestUseMiddlewareAbort(t *testing.T) {
	r := NewRouteGroup()
	r.Use(func(c *Context) {
		c.Redirect(302, "/login")
		c.Abort()
	})
	r.Get("/", func(c *Context) {
		c.String(200, "OK")
	})
	resp := processRequest(r, "GET", "/")
	assertEqual(t, 302, resp.StatusCode, "Status Code Error")
	assertEqual(t, "/login", resp.Header.Get("Location"), "Location header error")
}

func TestUseMiddlewareWithNext(t *testing.T) {
	r := NewRouteGroup()
	middle1Run := 0
	middle2Run := 0
	serveRun := 0
	r.Use(func(c *Context) {
		c.Next()
		middle1Run++
	})
	r.Use(func(c *Context) {
		c.Next()
		middle2Run++
	})
	r.Get("/", func(c *Context) {
		serveRun++
		c.String(200, "OK")
	})
	resp := processRequest(r, "GET", "/")
	assertEqual(t, 200, resp.StatusCode, "Status Code Error")
	assertEqual(t, 1, middle1Run, "Middle 1 run error")
	assertEqual(t, 1, middle2Run, "Middle 2 run error")
	assertEqual(t, 1, serveRun, "Serve run error")
}

func TestUseMiddlewareWithNextInAbort(t *testing.T) {
	r := NewRouteGroup()
	middle1Run := 0
	middle2Run := 0
	serveRun := 0
	r.Use(func(c *Context) {
		c.Next()
		middle1Run++
	})
	r.Use(func(c *Context) {
		c.Abort()
		middle2Run++
	})
	r.Get("/", func(c *Context) {
		serveRun++
		c.String(200, "OK")
	})
	resp := processRequest(r, "GET", "/")
	assertEqual(t, 200, resp.StatusCode, "Status Code Error")
	assertEqual(t, 1, middle1Run)
	assertEqual(t, 1, middle2Run)
	assertEqual(t, 0, serveRun)
}
