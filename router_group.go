package nazz

// 路由组
type RouterGroup struct {
	prefix         string
	beforeResponse []Middleware
	afterResponse  []Middleware
	server         *Server
}

func (this *Server) Group(prefix string, callback func(rg *RouterGroup)) {
	rg := &RouterGroup{
		server:         this,
		prefix:         prefix,
		beforeResponse: make([]Middleware, 0),
		afterResponse:  make([]Middleware, 0),
	}
	callback(rg)
}

func (this *RouterGroup) Use(middlewares ...string) *RouterGroup {
	for _, name := range middlewares {
		fn, ok := wares[name]
		if !ok {
			panic(name + " middleware not exist")
		}
		this.beforeResponse = append(this.beforeResponse, fn)
	}
	return this
}

func (this *RouterGroup) After(middlewares ...string) *RouterGroup {
	for _, name := range middlewares {
		fn, ok := wares[name]
		if !ok {
			panic(name + " middleware not exist")
		}
		this.afterResponse = append(this.afterResponse, fn)
	}
	return this
}

func (this *RouterGroup) GET(path string, handler Handler) Router {
	path = this.prefix + path
	r1 := staticRouter{
		Path:           path,
		Method:         "GET",
		Handler:        handler,
		BeforeResponse: this.beforeResponse,
		AfterResponse:  this.afterResponse,
	}
	if isStatic(path) {
		this.server.staticRouters["get:"+path] = &r1
		return &r1
	} else {
		prefix, re, params := parseDynamicRouter(path)
		r2 := &dynamicRouter{
			staticRouter: r1,
			re:           re,
			params:       params,
		}
		this.server.dynamicRouters["get:"+prefix] = r2
		return r2
	}
}

func (this *RouterGroup) POST(path string, handler Handler) Router {
	path = this.prefix + path
	r1 := staticRouter{
		Path:           path,
		Method:         "POST",
		Handler:        handler,
		BeforeResponse: this.beforeResponse,
		AfterResponse:  this.afterResponse,
	}
	if isStatic(path) {
		this.server.staticRouters["post:"+path] = &r1
		return &r1
	} else {
		prefix, re, params := parseDynamicRouter(path)
		r2 := &dynamicRouter{
			staticRouter: r1,
			re:           re,
			params:       params,
		}
		this.server.dynamicRouters["post:"+prefix] = r2
		return r2
	}
}
