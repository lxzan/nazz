package nazz

type Middleware func(ctx *Context) bool

var wares = map[string]Middleware{}

func Register(name string, fn Middleware) {
	wares[name] = fn
}
