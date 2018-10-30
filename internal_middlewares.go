package nazz

import "net/url"

// 解析query_string
func qsParser(ctx *Context) bool {
	if ctx.Request.URL.RawQuery == "" {
		return true
	}
	ctx.GET, _ = url.ParseQuery(ctx.Request.URL.RawQuery)
	return true
}
