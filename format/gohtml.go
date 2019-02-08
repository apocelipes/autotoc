package format

import "github.com/yosssi/gohtml"

func init() {
	SetFormatter("gohtml", goHtmlFormatter)
}

// 使用GoHTML格式化html
// 不支持自定义缩进
func goHtmlFormatter(html, _ string) string {
	return gohtml.Format(html)
}
