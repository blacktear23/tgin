package tgin

import (
	"log"
	"time"
)

func LoggerMiddleware(c *Context) {
	begin := time.Now()
	c.Next()
	ww, ok := c.Writer.(*ResponseWriterWrapper)
	if !ok {
		return
	}
	r := c.Request
	processTime := time.Now().Sub(begin).String()
	log.Printf("[Web] %d | %10s | %20s | %4s %s", ww.code, processTime, r.RemoteAddr, r.Method, r.URL.Path)
}
