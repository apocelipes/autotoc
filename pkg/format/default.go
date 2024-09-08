package format

import (
	"github.com/apocelipes/autotoc/pkg/format/defaultformatter"
)

func init() {
	SetFormatter("default", defaultHtmlFormatter)
}

// 使用defaultFormatter
// 尽量节省空间的使用，会将3个node合并在一行
func defaultHtmlFormatter(html, indent string) string {
	ret, _ := defaultformatter.FormatHtml(html, indent)
	return ret
}
