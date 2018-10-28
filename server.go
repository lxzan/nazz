package nazz

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type Handler func(ctx *Context) []byte

type Server struct {
	Port           int
	staticRouters  map[string]*staticRouter
	dynamicRouters map[string]*dynamicRouter
}

func NewServer() *Server {
	server := &Server{
		staticRouters:  map[string]*staticRouter{},
		dynamicRouters: map[string]*dynamicRouter{},
	}
	return server
}

func (this *Server) Get(path string, handler Handler) Router {
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

// 匹配动态路由
func (this *Server) matchDynamic(ctx *Context) (match bool, router *dynamicRouter) {
	paths := strings.Split(ctx.Request.URL.Path, "/")
	for i, _ := range paths {
		key := strings.ToLower(ctx.Request.Method) + ":" + strings.Join(paths[0:i], "/")
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
		ctx.PathParams[item.key] = paths[item.index]
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
		Response:   w,
		Request:    r,
		PathParams: Form{},
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
				for _, fn := range r1.BeforeResponse {
					if !fn(ctx) {
						return
					}
				}
				data := r2.Handler(ctx)
				ctx.Response.Write(data)
				for _, fn := range r1.AfterResponse {
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
