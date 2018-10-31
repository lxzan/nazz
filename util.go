package nazz

import (
	"regexp"
	"strconv"
	"strings"
	"time"
)

// 字符串模板
// Template("SELECT * FROM {{table}} WHERE id = {{id}}", "pre_forum_test", "1")
func Template(tpl string, args ...string) string {
	re := regexp.MustCompile(`{{.*?}}`)
	i := -1
	return re.ReplaceAllStringFunc(tpl, func(s string) string {
		i++
		return args[i]
	})
}

// 大驼峰转小驼峰
func toLowerCamel(s string) string {
	length := len(s)
	if length == 0 {
		return s
	}

	buf := []byte(s)
	if buf[0] >= 'A' && buf[0] <= 'Z' {
		buf[0] += 32
	}
	return string(buf)
}

// 大驼峰转小驼峰
func toUpperCamel(s string) string {
	length := len(s)
	if length == 0 {
		return s
	}

	buf := []byte(s)
	if buf[0] >= 'a' && buf[0] <= 'z' {
		buf[0] -= 32
	}
	return string(buf)
}

func ToString(num int64) string {
	return strconv.Itoa(int(num))
}

func Date(format string, timestamp ...int64) string {
	var ts = time.Now().Unix()
	if len(timestamp) > 0 {
		ts = timestamp[0]
	}
	var t = time.Unix(ts, 0)
	Y := t.Year()
	m := t.Month()
	d := t.Day()
	H := t.Hour()
	i := t.Minute()
	s := t.Second()
	format = strings.Replace(format, "Y", strconv.Itoa(Y), -1)
	format = strings.Replace(format, "m", strconv.Itoa(int(m)), -1)
	format = strings.Replace(format, "d", strconv.Itoa(d), -1)
	format = strings.Replace(format, "H", strconv.Itoa(H), -1)
	format = strings.Replace(format, "i", strconv.Itoa(i), -1)
	format = strings.Replace(format, "s", strconv.Itoa(s), -1)
	return format
}
