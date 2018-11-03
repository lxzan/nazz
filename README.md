# Nazz Framework
A high performance restful api framework

### Doc
- [Start](#start)
- [Router Group](#router-group)
- [Middleware](#middleware)

#### Start
```go
app := nazz.NewServer()
app.GET("/", func(ctx *nazz.Context) []byte {
    return ctx.JSON(nazz.J{
        "hello": "world!",
    })
})
app.Listen(8080)
```

#### Router Group
```go
app.Group("/home", func(rg *nazz.RouterGroup) {
    rg.GET("/", Demo)
    rg.POST("/test/:id", Test)
})

// get path param
func Test(ctx *nazz.Context) []byte {
    id := ctx.Request.Form.Get("id")
    println(id)
    return ctx.JSON(nazz.J{
        "success:": true,
    })
}
```

#### Middleware
```go
// You must register the middleware befor using it.
nazz.Register(nazz.BEFORE, "greet", func(ctx *nazz.Context) bool {
    println("Hello !")
    return true
})

// private use
app.GET("/", func(ctx *nazz.Context) []byte {
    return ctx.JSON(nazz.J{
        "hello": "world!",
    })
}).Use('greet')

// public use
app.Group("/home", func(rg *nazz.RouterGroup) {
    rg.Use("greet")
    rg.GET("/", Demo)
    rg.POST("/test/:id", Test)
})
```