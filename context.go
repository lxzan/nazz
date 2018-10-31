package nazz

import (
	"encoding/json"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
)

type Context struct {
	Response http.ResponseWriter
	Request  *http.Request
}

type J map[string]interface{}

const (
	urlEncode  = "application/x-www-form-urlencoded"
	jsonEncode = "application/json"
	formEncode = "multipart/form-data"
)

type Form map[string]string

func (this *Context) JSON(v interface{}, statusCode ...int) []byte {
	this.Response.Header().Set("Content-Type", "application/json; charset=UTF-8")
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

// parse url.Values to struct
func (this *Context) Parse(v url.Values, st interface{}) {
	val := reflect.ValueOf(st).Elem()
	t := reflect.TypeOf(st).Elem()
	n := t.NumField()
	for i := 0; i < n; i++ {
		tf := t.Field(i)
		name := toLowerCamel(tf.Name)
		tp := tf.Type.String()

		switch tp {
		case "string":
			val.Field(i).SetString(v.Get(name))
			break
		case "int":
			num, _ := strconv.Atoi(v.Get(name))
			val.Field(i).SetInt(int64(num))
			break
		case "int64":
			num, _ := strconv.Atoi(v.Get(name))
			val.Field(i).SetInt(int64(num))
			break
		case "[]string":
			val.Field(i).Set(reflect.ValueOf(v[name+"[]"]))
			break
		case "[]int":
			var arr = make([]int, 0)
			ss := v[name+"[]"]
			for _, s := range ss {
				num, _ := strconv.Atoi(s)
				arr = append(arr, num)
			}
			val.Field(i).Set(reflect.ValueOf(arr))
			break
		case "[]int64":
			var arr = make([]int64, 0)
			ss := v[name+"[]"]
			for _, s := range ss {
				num, _ := strconv.Atoi(s)
				arr = append(arr, int64(num))
			}
			val.Field(i).Set(reflect.ValueOf(arr))
		}
	}
}

func (this *Context) ParseForm(st interface{}) {
	this.Parse(this.Request.PostForm, st)
}

func (this *Context) ParseQuery(st interface{}) {
	this.Parse(this.Request.Form, st)
}

func (this *Context) GetBody(input interface{}) []byte {
	var body = make([]byte, this.Request.ContentLength)
	this.Request.Body.Read(body)
	return body
}
