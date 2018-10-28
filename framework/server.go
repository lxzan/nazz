package nazz

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type Handler func(ctx *Context)

type Server struct {
	Port           int
	staticRouters  map[string]*staticRouter
	dynamicRouters map[string]*dynamicRouter
	//MiddleWares
}

func NewServer() *Server {
	server := &Server{
		staticRouters:  make(map[string]*staticRouter),
		dynamicRouters: make(map[string]*dynamicRouter),
	}
	return server
}

func (this *Server) Get(path string, handler Handler) {
	r := staticRouter{
		Path:    path,
		Method:  "GET",
		Handler: handler,
	}
	if isStatic(path) {
		this.staticRouters["get:"+path] = &r
	} else {
		prefix, re, params := parseDynamicRouter(path)
		this.dynamicRouters["get:"+prefix] = &dynamicRouter{
			staticRouter: r,
			re:           re,
			params:       params,
		}
	}
}

// 匹配动态路由
func (this *Server) matchDynamic(ctx *Context) (match bool, router *dynamicRouter) {
	paths := strings.Split(ctx.HttpRequest.URL.Path, "/")
	for i, _ := range paths {
		key := strings.ToLower(ctx.HttpRequest.Method) + ":" + strings.Join(paths[0:i], "/")
		r, ok := this.dynamicRouters[key]
		if ok && r.re.MatchString(ctx.HttpRequest.URL.Path) {
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
	Writer      http.ResponseWriter
	HttpRequest *http.Request
	Callback    func(ctx *Context)
}

func NewHandler(callback func(ctx *Context)) *globalHandler {
	obj := &globalHandler{}
	obj.Callback = callback
	return obj
}

func (this *globalHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	this.Callback(&Context{
		Writer:      w,
		HttpRequest: r,
		PathParams:  Form{},
	})
}

func (this *Server) Listen(port int) {
	addr := ":" + strconv.Itoa(port)
	println(fmt.Sprintf("Nazz is listening on port %d", port))

	http.ListenAndServe(addr, NewHandler(func(ctx *Context) {
		path := ctx.HttpRequest.URL.Path
		key := strings.ToLower(ctx.HttpRequest.Method) + ":" + path
		r1, ok := this.staticRouters[key]
		isMatch := false
		if ok {
			isMatch = true
			r1.Handler(ctx)
		} else {
			m, r2 := this.matchDynamic(ctx)
			if m {
				isMatch = true
				r2.Handler(ctx)
			}
		}

		if !isMatch {
			ctx.Render([]byte("<h1>404 Not Found</h1>"), 404)
		}
	}))
}
