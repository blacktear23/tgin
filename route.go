package tgin

import (
	"log"
	"net/http"
	"time"
)

var (
	_ http.Handler = (*RouteGroup)(nil)
)

type RouteHandler func(c *Context)

type RouteGroup struct {
	prefix      string
	mux         *http.ServeMux
	middlewares []RouteHandler
}

func NewRouteGroup() *RouteGroup {
	return &RouteGroup{
		prefix:      "",
		mux:         http.NewServeMux(),
		middlewares: []RouteHandler{},
	}
}

func (rg *RouteGroup) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	begin := time.Now()
	ww := &ResponseWriterWrapper{
		writer: w,
		code:   200,
	}
	ctx := newContext(ww, r)
	ww.ctx = ctx
	returned := false
	for _, middleware := range rg.middlewares {
		middleware(ctx)
		if ctx.aborted {
			returned = true
			break
		}
	}
	if !returned {
		rg.mux.ServeHTTP(ww, r)
	}
	processTime := time.Now().Sub(begin).String()
	log.Printf("[Web] %d | %10s | %20s | %4s %s", ww.code, processTime, r.RemoteAddr, r.Method, r.URL.Path)
}

func (rg *RouteGroup) getPath(path string) string {
	return rg.prefix + path
}

func (rg *RouteGroup) getContext(w http.ResponseWriter, r *http.Request) *Context {
	if rww, ok := w.(*ResponseWriterWrapper); ok {
		if rww.ctx != nil {
			return rww.ctx
		}
	}
	return newContext(w, r)
}

func (rg *RouteGroup) handle(method, path string, handler RouteHandler) {
	rg.mux.HandleFunc(rg.getPath(path), func(w http.ResponseWriter, r *http.Request) {
		ctx := rg.getContext(w, r)
		if method != ctx.Method {
			ctx.Text(404, "404 page not found\n")
			return
		}
		handler(ctx)
	})
}

func (rg *RouteGroup) Group(prefix string) *RouteGroup {
	return &RouteGroup{
		prefix: rg.getPath(prefix),
		mux:    rg.mux,
	}
}

func (rg *RouteGroup) Any(path string, handler RouteHandler) {
	rg.mux.HandleFunc(rg.getPath(path), func(w http.ResponseWriter, r *http.Request) {
		ctx := rg.getContext(w, r)
		handler(ctx)
	})
}

func (rg *RouteGroup) Get(path string, handler RouteHandler) {
	rg.handle("GET", path, handler)
}

func (rg *RouteGroup) GET(path string, handler RouteHandler) {
	rg.Get(path, handler)
}

func (rg *RouteGroup) Post(path string, handler RouteHandler) {
	rg.handle("POST", path, handler)
}

func (rg *RouteGroup) POST(path string, handler RouteHandler) {
	rg.Post(path, handler)
}

func (rg *RouteGroup) Put(path string, handler RouteHandler) {
	rg.handle("PUT", path, handler)
}

func (rg *RouteGroup) PUT(path string, handler RouteHandler) {
	rg.Put(path, handler)
}

func (rg *RouteGroup) Delete(path string, handler RouteHandler) {
	rg.handle("DELETE", path, handler)
}

func (rg *RouteGroup) DELETE(path string, handler RouteHandler) {
	rg.Delete(path, handler)
}

func (rg *RouteGroup) Head(path string, handler RouteHandler) {
	rg.handle("HEAD", path, handler)
}

func (rg *RouteGroup) HEAD(path string, handler RouteHandler) {
	rg.Head(path, handler)
}

func (rg *RouteGroup) Options(path string, handler RouteHandler) {
	rg.handle("OPTIONS", path, handler)
}

func (rg *RouteGroup) OPTIONS(path string, handler RouteHandler) {
	rg.Options(path, handler)
}

func (rg *RouteGroup) StaticFile(path, filePath string) {
	handler := func(c *Context) {
		c.File(filePath)
	}
	rg.Get(path, handler)
	rg.Head(path, handler)
}

func (rg *RouteGroup) Use(middlewares ...RouteHandler) {
	rg.middlewares = append(rg.middlewares, middlewares...)
}