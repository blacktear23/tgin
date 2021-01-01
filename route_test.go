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
	assertBody(t, resp, "Hello World", "Body not correct")

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
	middle3Run := 0
	serveRun := 0
	r.Use(func(c *Context) {
		c.Next()
		middle1Run++
	})
	r.Use(func(c *Context) {
		c.Next()
		middle2Run++
	})
	r.Use(func(c *Context) {
		middle3Run++
	})
	r.Get("/", func(c *Context) {
		serveRun++
		c.String(200, "OK")
	})
	resp := processRequest(r, "GET", "/")
	assertEqual(t, 200, resp.StatusCode, "Status Code Error")
	assertEqual(t, 1, middle1Run, "Middle 1 run error")
	assertEqual(t, 1, middle2Run, "Middle 2 run error")
	assertEqual(t, 1, middle3Run, "Middle 2 run error")
	assertEqual(t, 1, serveRun, "Serve run error")
}

func TestUseMiddlewareWithNextInAbort(t *testing.T) {
	r := NewRouteGroup()
	middle1Run := 0
	middle2Run := 0
	middle3Run := 0
	serveRun := 0
	r.Use(func(c *Context) {
		c.Next()
		middle1Run++
	})
	r.Use(func(c *Context) {
		c.Abort()
		middle2Run++
	})
	r.Use(func(c *Context) {
		middle3Run++
	})
	r.Get("/", func(c *Context) {
		serveRun++
		c.String(200, "OK")
	})
	resp := processRequest(r, "GET", "/")
	assertEqual(t, 200, resp.StatusCode, "Status Code Error")
	assertEqual(t, 1, middle1Run)
	assertEqual(t, 1, middle2Run)
	assertEqual(t, 0, middle3Run)
	assertEqual(t, 0, serveRun)
}

func TestSamePathWithDifferentMethodHandlers(t *testing.T) {
	r := NewRouteGroup()
	r.Get("/hello", func(c *Context) {
		c.String(200, "Get Hello World")
	})
	r.Post("/hello", func(c *Context) {
		c.String(200, "Post Hello World")
	})
	r.Put("/hello", func(c *Context) {
		c.String(200, "Put Hello World")
	})

	resp := processRequest(r, "GET", "/hello")
	assertEqual(t, 200, resp.StatusCode)
	assertBody(t, resp, "Get Hello World")

	resp = processRequest(r, "POST", "/hello")
	assertEqual(t, 200, resp.StatusCode)
	assertBody(t, resp, "Post Hello World")

	resp = processRequest(r, "PUT", "/hello")
	assertEqual(t, 200, resp.StatusCode)
	assertBody(t, resp, "Put Hello World")

	resp = processRequest(r, "DELETE", "/hello")
	assertEqual(t, 404, resp.StatusCode)
}

func TestMiddlewareInDifferentGroups(t *testing.T) {
	r := NewRouteGroup()
	middle1Run := 0
	middle2Run := 0
	r.Use(func(c *Context) {
		middle1Run++
	})
	nr := r.Group("/test")
	nr.Use(func(c *Context) {
		middle2Run++
	})

	r.Get("/hello", func(c *Context) {
		c.String(200, "Hello World")
	})
	nr.Get("/hello", func(c *Context) {
		c.String(200, "Hello World")
	})

	resp := processRequest(r, "GET", "/hello")
	assertEqual(t, 200, resp.StatusCode)
	assertEqual(t, 1, middle1Run)
	assertEqual(t, 0, middle2Run)

	resp = processRequest(r, "GET", "/test/hello")
	assertEqual(t, 200, resp.StatusCode)
	assertEqual(t, 2, middle1Run)
	assertEqual(t, 1, middle2Run)
}

func TestMiddlewareWithAbortBefore404(t *testing.T) {
	r := NewRouteGroup()
	middle1Run := 0
	handlerRun := 0
	r.Use(func(c *Context) {
		middle1Run++
		c.String(500, "Nothing")
		c.Abort()
	})

	r.Get("/hello", func(c *Context) {
		handlerRun++
		c.String(200, "Hello World")
	})

	resp := processRequest(r, "PUT", "/hello")
	assertEqual(t, 500, resp.StatusCode)
	assertEqual(t, 1, middle1Run)
	assertEqual(t, 0, handlerRun)
}

func TestMiddlewareWithAny(t *testing.T) {
	r := NewRouteGroup()
	middle1Run := 0
	handlerRun := 0
	r.Use(func(c *Context) {
		middle1Run++
	})

	r.Any("/hello", func(c *Context) {
		handlerRun++
		c.String(200, "Hello World")
	})

	resp := processRequest(r, "PUT", "/hello")
	assertEqual(t, 200, resp.StatusCode)
	assertEqual(t, 1, middle1Run)
	assertEqual(t, 1, handlerRun)
}

func TestMiddlewareAbortWithAny(t *testing.T) {
	r := NewRouteGroup()
	middle1Run := 0
	handlerRun := 0
	r.Use(func(c *Context) {
		middle1Run++
		c.String(500, "Nothing")
		c.Abort()
	})

	r.Any("/hello", func(c *Context) {
		handlerRun++
		c.String(200, "Hello World")
	})

	resp := processRequest(r, "PUT", "/hello")
	assertEqual(t, 500, resp.StatusCode)
	assertEqual(t, 1, middle1Run)
	assertEqual(t, 0, handlerRun)
}

func TestRecoveryMiddleware(t *testing.T) {
	r := NewRouteGroup()
	middle1Run := 0
	r.Use(RecoveryMiddleware)
	r.Use(func(c *Context) {
		middle1Run++
	})
	r.Get("/", func(c *Context) {
		panic("this is a test panic")
	})

	resp := processRequest(r, "GET", "/")
	assertEqual(t, 500, resp.StatusCode)
}

func TestRecoveryMiddlewareWithPanicInMiddleware(t *testing.T) {
	r := NewRouteGroup()
	middle1Run := 0
	middle2Run := 0
	middle3Run := 0
	handlerRun := 0
	r.Use(RecoveryMiddleware)
	r.Use(func(c *Context) {
		middle1Run++
	})
	r.Use(func(c *Context) {
		middle2Run++
		panic("this is a test panic from middleware")
	})
	r.Use(func(c *Context) {
		middle3Run++
	})
	r.Get("/", func(c *Context) {
		handlerRun++
		c.String(200, "Hello")
	})

	resp := processRequest(r, "GET", "/")
	assertEqual(t, 500, resp.StatusCode)
	assertEqual(t, 1, middle1Run)
	assertEqual(t, 1, middle2Run)
	assertEqual(t, 0, middle3Run)
	assertEqual(t, 0, handlerRun)
}

func TestStaticFileMiddleware(t *testing.T) {
	r := NewRouteGroup()
	r.UseGlobal(StaticFileMiddleware("/", "/tmp", true))
	r.Get("/hello", func(c *Context) {
		c.String(200, "Hello world")
	})
	resp := processRequest(r, "GET", "/")
	assertEqual(t, 200, resp.StatusCode)
	resp = processRequest(r, "GET", "/hello")
	assertEqual(t, 200, resp.StatusCode)
	assertBody(t, resp, "Hello world")
}
