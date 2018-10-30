package nazz

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"net/url"
)

type Handler func(ctx *Context) []byte

type Server struct {
	staticRouters  map[string]*staticRouter
	dynamicRouters map[string]*dynamicRouter
}

func NewServer() *Server {
	server := &Server{
		staticRouters:  map[string]*staticRouter{},
		dynamicRouters: map[string]*dynamicRouter{},
	}

	Register(GLOBAL_BEFORE, "qs_parser", qsParser)
	return server
}

func (this *Server) GET(path string, handler Handler) Router {
	r1 := staticRouter{
		Path:    path,
		Method:  "GET",
		Handler: handler,
	}
	if isStatic(path) {
		this.staticRouters["get:"+path] = &r1
		return &r1
	} else {
		prefix, re, params := parseDynamicRouter(path)
		r2 := &dynamicRouter{
			staticRouter: r1,
			re:           re,
			params:       params,
		}
		this.dynamicRouters["get:"+prefix] = r2
		return r2
	}
}

func (this *Server) POST(path string, handler Handler) Router {
	r1 := staticRouter{
		Path:    path,
		Method:  "POST",
		Handler: handler,
	}
	if isStatic(path) {
		this.staticRouters["post:"+path] = &r1
		return &r1
	} else {
		prefix, re, params := parseDynamicRouter(path)
		r2 := &dynamicRouter{
			staticRouter: r1,
			re:           re,
			params:       params,
		}
		this.dynamicRouters["post:"+prefix] = r2
		return r2
	}
}

// 匹配动态路由
func (this *Server) matchDynamic(ctx *Context) (match bool, router *dynamicRouter) {
	paths := strings.Split(ctx.Request.URL.Path, "/")
	for i, _ := range paths {
		key := strings.ToLower(ctx.Request.Method) + ":" + strings.Join(paths[0:i+1], "/")
		r, ok := this.dynamicRouters[key]
		if ok && r.re.MatchString(ctx.Request.URL.Path) {
			router = r
			match = true
			break
		}
	}

	if !match {
		return false, nil
	}
	for _, item := range router.params {
		ctx.PATHINFO[item.key] = paths[item.index]
	}

	for k, v := range ctx.PATHINFO {
		ctx.GET.Add(k, v)
	}
	return true, router
}

type globalHandler struct {
	Callback func(ctx *Context)
}

func NewHandler(callback func(ctx *Context)) *globalHandler {
	obj := &globalHandler{}
	obj.Callback = callback
	return obj
}

func (this *globalHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	this.Callback(&Context{
		Response: w,
		Request:  r,
		PATHINFO: Form{},
		GET:      url.Values{},
	})
}

func (this *Server) Listen(port int) {
	addr := ":" + strconv.Itoa(port)
	println(fmt.Sprintf("Nazz is listening on port %d", port))

	http.ListenAndServe(addr, NewHandler(func(ctx *Context) {
		for _, fn := range globalBeforeWares {
			if !fn(ctx) {
				return
			}
		}

		ctx.Request.URL.Path = filterLastSlash(ctx.Request.URL.Path)
		path := ctx.Request.URL.Path
		key := strings.ToLower(ctx.Request.Method) + ":" + path
		r1, ok := this.staticRouters[key]
		isMatch := false
		if ok {
			isMatch = true
			for _, fn := range r1.BeforeResponse {
				if !fn(ctx) {
					return
				}
			}
			data := r1.Handler(ctx)
			ctx.Response.Write(data)
			for _, fn := range r1.AfterResponse {
				if !fn(ctx) {
					return
				}
			}
		} else {
			m, r2 := this.matchDynamic(ctx)
			if m {
				isMatch = true
				for _, fn := range r2.BeforeResponse {
					if !fn(ctx) {
						return
					}
				}
				data := r2.Handler(ctx)
				ctx.Response.Write(data)
				for _, fn := range r2.AfterResponse {
					if !fn(ctx) {
						return
					}
				}
			}
		}

		if !isMatch {
			ctx.Render([]byte("<h1>404 Not Found</h1>"), 404)
		}

		for _, fn := range globalAfterWares {
			if !fn(ctx) {
				return
			}
		}
	}))
}

// 过滤最后的斜杠
func filterLastSlash(path string) string {
	length := len(path)
	if length == 0 || length == 1 || path[length-1] != '/' {
		return path
	}
	return string(path[0 : length-1])
}
