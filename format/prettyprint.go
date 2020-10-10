package format

import "github.com/hermanschaaf/prettyprint"

func init() {
	SetFormatter("prettyprint", prettyPrintFormatter)
}

// prettyPrintFormatter 使用prettyprint格式化html
func prettyPrintFormatter(html, indent string) string {
	ret, _ := prettyprint.Prettify(html, indent)
	return ret
}
