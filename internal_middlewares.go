package nazz

import (
	"log"
	"mime"
	"mime/multipart"
	"net/url"
	"os"
)

// 解析query_string和文件
func paramParser(ctx *Context) bool {
	if ctx.Request.URL.RawQuery != "" {
		ctx.Request.Form, _ = url.ParseQuery(ctx.Request.URL.RawQuery)
	}

	contentType := ctx.Request.Header.Get("Content-Type")
	mediaType, params, _ := mime.ParseMediaType(contentType)

	if mediaType == urlEncode {
		var buf = make([]byte, ctx.Request.ContentLength)
		ctx.Request.Body.Read(buf)
		ctx.Request.PostForm, _ = url.ParseQuery(string(buf))
	}
	if mediaType == formEncode {
		boundary, _ := params["boundary"]
		reader := multipart.NewReader(ctx.Request.Body, boundary)
		ctx.Request.MultipartForm, _ = reader.ReadForm(ctx.Request.ContentLength)
	}
	return true
}

var logger = log.New(os.Stdout, "Request: ", 3)

func accessLog(ctx *Context) bool {
	logger.Printf("IP=%s, Method=%s, URI=%s, UA=%s", ctx.Request.RemoteAddr, ctx.Request.Method, ctx.Request.URL.RequestURI(), ctx.Request.Header.Get("User-Agent"))
	return true
}
