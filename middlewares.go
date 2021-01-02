package tgin

import (
	"strings"
)

type RouteHandlerChain []RouteHandler

type MiddlewareTree struct {
	tree map[string]RouteHandlerChain
}

func newMiddlewareTree() *MiddlewareTree {
	return &MiddlewareTree{
		tree: make(map[string]RouteHandlerChain),
	}
}

func (mt *MiddlewareTree) Add(prefix string, handler RouteHandler) {
	if prefix == "" {
		prefix = "/"
	} else {
		prefix = strings.TrimSuffix(prefix, "/")
	}
	rhc, have := mt.tree[prefix]
	if have {
		mt.tree[prefix] = append(rhc, handler)
	} else {
		mt.tree[prefix] = RouteHandlerChain{handler}
	}
}

func (mt *MiddlewareTree) BuildMiddlewares(path string) RouteHandlerChain {
	ret := RouteHandlerChain{}
	parts := strings.Split(path, "/")
	cp := ""
	for i, p := range parts {
		if i == 0 {
			rhc, have := mt.tree["/"]
			if have {
				ret = appendRouterHandlerChain(ret, rhc)
			}
		} else if p != "" {
			cp += ("/" + p)
			rhc, have := mt.tree[cp]
			if have {
				ret = appendRouterHandlerChain(ret, rhc)
			}
		}
	}
	return ret
}

func appendRouterHandlerChain(current, appended RouteHandlerChain) RouteHandlerChain {
	for _, rh := range appended {
		current = append(current, rh)
	}
	return current
}
