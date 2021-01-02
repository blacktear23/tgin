package tgin

import (
	"testing"
)

func TestBuildMiddlewares(t *testing.T) {
	mt := newMiddlewareTree()
	mt.Add("", func(c *Context) {})
	mt.Add("", func(c *Context) {})
	mt.Add("/test", func(c *Context) {})
	mt.Add("/test/demo", func(c *Context) {})
	mt.Add("/test/demo", func(c *Context) {})
	mt.Add("/test/shake", func(c *Context) {})
	mt.Add("/make/demo", func(c *Context) {})

	assertEqual(t, 2, len(mt.BuildMiddlewares("/")))
	assertEqual(t, 2, len(mt.BuildMiddlewares("/asdf")))
	assertEqual(t, 2, len(mt.BuildMiddlewares("/asdf/")))

	assertEqual(t, 3, len(mt.BuildMiddlewares("/test")))
	assertEqual(t, 3, len(mt.BuildMiddlewares("/test/")))
	assertEqual(t, 3, len(mt.BuildMiddlewares("/test/api")))

	assertEqual(t, 5, len(mt.BuildMiddlewares("/test/demo")))
	assertEqual(t, 5, len(mt.BuildMiddlewares("/test/demo/")))
	assertEqual(t, 5, len(mt.BuildMiddlewares("/test/demo/api")))

	assertEqual(t, 4, len(mt.BuildMiddlewares("/test/shake")))
	assertEqual(t, 4, len(mt.BuildMiddlewares("/test/shake/")))
	assertEqual(t, 4, len(mt.BuildMiddlewares("/test/shake/api")))

	assertEqual(t, 2, len(mt.BuildMiddlewares("/make")))
	assertEqual(t, 2, len(mt.BuildMiddlewares("/make/")))
	assertEqual(t, 2, len(mt.BuildMiddlewares("/make/api")))

	assertEqual(t, 3, len(mt.BuildMiddlewares("/make/demo")))
	assertEqual(t, 3, len(mt.BuildMiddlewares("/make/demo/")))
	assertEqual(t, 3, len(mt.BuildMiddlewares("/make/demo/api")))
}
