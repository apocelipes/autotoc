package parser

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
}
