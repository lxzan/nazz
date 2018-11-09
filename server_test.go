package nazz

import "testing"

func TestNewServer(t *testing.T) {
	app := NewServer()
	app.GET("/", func(ctx *Context) []byte {
		return ctx.JSON(J{
			"hello": "<lxz>",
		})
	})
	app.Listen(8081)
}
