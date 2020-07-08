package tgin

import (
	"net/http"
)

type Engine struct {
	*RouteGroup
}

func New() *Engine {
	return &Engine{
		RouteGroup: NewRouteGroup(),
	}
}

func (e *Engine) Run(addr string) error {
	server := &http.Server{
		Addr:           addr,
		Handler:        e.RouteGroup,
		MaxHeaderBytes: 1 << 20,
	}
	return server.ListenAndServe()
}

func (e *Engine) RunTLS(addr, certFile, keyFile string) error {
	server := &http.Server{
		Addr:           addr,
		Handler:        e.RouteGroup,
		MaxHeaderBytes: 1 << 20,
	}
	return server.ListenAndServeTLS(certFile, keyFile)
}
