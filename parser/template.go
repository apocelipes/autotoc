package parser

const (
	// 不包含子目录的目录模板
	catalogItemTemplate = "<li><a href=\"#%s\">%s</a></li>\n"
	// 包含子目录的目录模板
	catalogItemWithChildren = "<li>\n<a href=\"#%s\">%s</a>\n"
)
