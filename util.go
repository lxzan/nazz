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
