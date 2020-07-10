package tgin

import (
	"testing"
)

func TestDefaultEngine(t *testing.T) {
	e := Default()
	e.Get("/", func(c *Context) {
		c.String(200, "OK")
	})
	resp := processRequest(e.RouteGroup, "GET", "/")
	assertEqual(t, 200, resp.StatusCode)
}

func TestDefaultEngineWithPanic(t *testing.T) {
	e := Default()
	e.Get("/", func(c *Context) {
		panic("this is a test panic")
	})
	resp := processRequest(e.RouteGroup, "GET", "/")
	assertEqual(t, 500, resp.StatusCode)
}
