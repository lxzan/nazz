package nazz

import (
	"encoding/json"
	"net/http"
)

type Context struct {
	Response   http.ResponseWriter
	Request    *http.Request
	PathParams Form
}

type J map[string]interface{}

type Form map[string]string

func (this *Context) JSON(v interface{}, statusCode ...int) []byte {
	if len(statusCode) > 0 {
		this.Response.WriteHeader(statusCode[0])
	}

	data, _ := json.Marshal(v)
	return data
}

func (this *Context) Render(data []byte, statusCode ...int) error {
	if len(statusCode) > 0 {
		this.Response.WriteHeader(statusCode[0])
	}
	_, err := this.Response.Write(data)
	return err
}

func (this *Context) SetHeaders(headers Form) {
	for k, v := range headers {
		this.Response.Header().Set(k, v)
	}
}
