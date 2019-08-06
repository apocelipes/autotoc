package parser

// 对URL进行编码的函数，传入URL，传出编码后的结果
type ContentEncoder func(string) string

// 用于向parseMarkdown的参数传递
type ParseOption struct {
	// 开始解析的标题层次
	TopTag string
	// 待扫描的标题行的类型
	ScanType string
	// 从tocMark后的位置开始解析标题
	TocMark string
	// 根据规则或标题内容过滤title
	Filter TitleFilter
	// 对标题的URL进行encoding
	Encoder ContentEncoder
}
