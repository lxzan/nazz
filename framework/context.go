package nazz

import (
	"encoding/json"
	"net/http"
)

type Context struct {
	Writer      http.ResponseWriter
	HttpRequest *http.Request
	PathParams  Form
}

type J map[string]interface{}
type Form map[string]string

func (this *Context) JSON(v interface{}, statusCode ...int) error {
	if len(statusCode) > 0 {
		this.Writer.WriteHeader(statusCode[0])
	}

	data, err1 := json.Marshal(v)
	if err1 != nil {
		return err1
	}

	_, err2 := this.Writer.Write(data)
	if err2 != nil {
		return err2
	}
	return nil
}

func (this *Context) Render(data []byte, statusCode ...int) error {
	if len(statusCode) > 0 {
		this.Writer.WriteHeader(statusCode[0])
	}
	_, err := this.Writer.Write(data)
	return err
}

func (this *Context) SetHeaders(headers Form) {
	for k, v := range headers {
		this.Writer.Header().Set(k, v)
	}
}
