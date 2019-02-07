package format

import "github.com/yosssi/gohtml"

func init() {
	SetFormatter("gohtml", goHtmlFormatter)
}

// 使用GoHTML格式化html
func goHtmlFormatter(html string) string {
	return gohtml.Format(html)
}
