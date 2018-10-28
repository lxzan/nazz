package nazz

import (
	"regexp"
	"strings"
)

type Router interface {
	Use(middlewares ...string) Router
	After(middlewares ...string) Router
}

// 静态路由
type staticRouter struct {
	Path           string
	Method         string
	Handler        Handler
	BeforeResponse []Middleware
	AfterResponse  []Middleware
}

// 请求前中间件
func (this *staticRouter) Use(middlewares ...string) Router {
	for _, name := range middlewares {
		fn, ok := wares[name]
		if !ok {
			panic(name + " middleware not exist")
		}
		this.BeforeResponse = append(this.BeforeResponse, fn)
	}
	return this
}

func (this *staticRouter) After(middlewares ...string) Router {
	for _, name := range middlewares {
		fn, ok := wares[name]
		if !ok {
			panic(name + " middleware not exist")
		}
		this.AfterResponse = append(this.AfterResponse, fn)
	}
	return this
}

type dynamicRouterParam struct {
	key   string
	index int
}

type dynamicRouter struct {
	staticRouter
	re     *regexp.Regexp
	params []dynamicRouterParam
}

// 是否为静态路由
func isStatic(path string) bool {
	arr := strings.Split(path, "/")
	for _, item := range arr {
		if item == "" {
			continue
		}
		if item[0] == ':' {
			return false
		}
	}
	return true
}

// 解析动态路由
func parseDynamicRouter(path string) (prefix string, re *regexp.Regexp, params []dynamicRouterParam) {
	params = make([]dynamicRouterParam, 0)
	flag := true
	patterns := make([]string, 0)
	paths := make([]string, 0)

	arr := strings.Split(path, "/")
	for i, item := range arr {
		var length = len(item)
		pattern := item
		if length > 0 && item[0] == ':' {
			flag = false
			pattern = "[0-9a-zA-z]+"
			params = append(params, dynamicRouterParam{
				key:   item[1:length],
				index: i,
			})
		}
		if flag {
			paths = append(paths, item)
		}
		patterns = append(patterns, pattern)
	}

	prefix = strings.Join(paths, "/")
	re = regexp.MustCompile("^" + strings.Join(patterns, `\/`) + "$")
	return prefix, re, params
}
