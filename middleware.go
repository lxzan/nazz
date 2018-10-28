package nazz

type Middleware func(ctx *Context) bool

type MiddlewareType int

const (
	BEFORE MiddlewareType = iota
	AFTER
	GLOBAL_BEFORE
	GLOBAL_AFTER
)

var wares = map[string]Middleware{}

var globalBeforeWares = make([]Middleware, 0)

var globalAfterWares = make([]Middleware, 0)

func Register(middlewareType MiddlewareType, name string, fn Middleware) {
	switch middlewareType {
	case BEFORE:
		wares[name] = fn
		break
	case AFTER:
		wares[name] = fn
		break
	case GLOBAL_BEFORE:
		globalBeforeWares = append(globalBeforeWares, fn)
		break
	case GLOBAL_AFTER:
		globalAfterWares = append(globalAfterWares, fn)
		break
	}
}
