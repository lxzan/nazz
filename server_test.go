package nazz

import "testing"

func TestNewServer(t *testing.T) {
	app := NewServer()
	app.GET("/", func(ctx *Context) []byte {
		return ctx.Redirect("https://github.com/boltdb/bolt#read-write-transactions")
	})
	app.Listen(8081)
}
