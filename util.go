package nazz

import (
	"regexp"
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
