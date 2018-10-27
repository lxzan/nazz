package nazz

import (
	"regexp"
	"strings"
)

// 是否为静态路由
func isStatic(path string) bool {
	arr := strings.Split(path, "/")
	for _, item := range arr {
		if item == "" {
			continue
		}
		if item[0] == ':' {
			return false
		}
	}
	return true
}

// 解析动态路由
func parseDynamicRouter(path string) (prefix string, re *regexp.Regexp, params []dynamicRouterParam) {
	params = make([]dynamicRouterParam, 0)
	flag := true
	patterns := make([]string, 0)
	paths := make([]string, 0)

	arr := strings.Split(path, "/")
	for i, item := range arr {
		var length = len(item)
		pattern := item
		if length > 0 && item[0] == ':' {
			flag = false
			pattern = "[0-9a-zA-z]+"
			params = append(params, dynamicRouterParam{
				key:   item[1 : length],
				index: i,
			})
		}
		if flag {
			paths = append(paths, item)
		}
		patterns = append(patterns, pattern)
	}

	prefix = strings.Join(paths, "/")
	re = regexp.MustCompile(strings.Join(patterns, `\/`) + "$")
	return prefix, re, params
}

// 字符串模板
// Template("SELECT * FROM {table} WHERE id = {id}", "pre_forum_test", "1")
func Template(tpl string, args ...string) string {
	re := regexp.MustCompile(`{.*?}`)
	i := -1
	return re.ReplaceAllStringFunc(tpl, func(s string) string {
		i++
		return args[i]
	})
}
