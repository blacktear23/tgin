package tgin

import (
	"net/http"
)

var (
	_ http.Handler = (*RouteGroup)(nil)
)

type RouteHandler func(c *Context)

type handlerFunctions map[string]RouteHandler

type RouteGroup struct {
	prefix      string
	handlers    map[string]handlerFunctions
	mux         *http.ServeMux
	middlewares []RouteHandler
}

func NewRouteGroup() *RouteGroup {
	return &RouteGroup{
		prefix:      "",
		mux:         http.NewServeMux(),
		middlewares: []RouteHandler{},
		handlers:    map[string]handlerFunctions{},
	}
}

func (rg *RouteGroup) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ww := &ResponseWriterWrapper{
		ResponseWriter: w,
		code:           200,
	}
	if hj, ok := w.(http.Hijacker); ok {
		ww.Hijacker = hj
	}
	ctx := newContext(ww, r)
	ww.ctx = ctx
	ctx.mux = rg.mux
	rg.mux.ServeHTTP(ww, r)
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
	fullPath := rg.getPath(path)
	hfs, have := rg.handlers[fullPath]
	if have {
		hfs[method] = handler
	} else {
		nhfs := handlerFunctions{}
		nhfs[method] = handler
		rg.handlers[fullPath] = nhfs
		rg.mux.HandleFunc(fullPath, func(w http.ResponseWriter, r *http.Request) {
			ctx := rg.getContext(w, r)
			ctx.middlewares = rg.middlewares
			ctx.handler = func(c *Context) {
				lhfs, have := rg.handlers[fullPath]
				if !have {
					ctx.Text(404, "404 page not found\n")
					return
				}
				hdl, have := lhfs[ctx.Method]
				if !have {
					ctx.Text(404, "404 page not found\n")
					return
				}
				hdl(c)
			}
			ctx.Next()
		})
	}
}

func (rg *RouteGroup) Group(prefix string) *RouteGroup {
	newMiddlewares := make([]RouteHandler, len(rg.middlewares))
	for i, m := range rg.middlewares {
		newMiddlewares[i] = m
	}
	return &RouteGroup{
		prefix:      rg.getPath(prefix),
		mux:         rg.mux,
		middlewares: newMiddlewares,
		handlers:    map[string]handlerFunctions{},
	}
}

func (rg *RouteGroup) Any(path string, handler RouteHandler) {
	rg.mux.HandleFunc(rg.getPath(path), func(w http.ResponseWriter, r *http.Request) {
		ctx := rg.getContext(w, r)
		ctx.middlewares = rg.middlewares
		ctx.handler = handler
		ctx.Next()
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
