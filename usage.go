package main

import (
	"fmt"
	"os"
)

// 参数选项说明
const (
	topTagUsage            = "设置作为目录顶层项的tag，将从指定tag开始解析标题"
	formatterUsage         = "选择格式化html代码的方式，目前只支持default和prettyprint(output为markdown时不支持)"
	catalogIdUsage         = "目录的html id(output为markdown时不支持)"
	catalogTitleUsage      = "目录的标题"
	catalogOutputTypeUsage = "输出的目录格式，可以为html或md(markdown)"
	catalogScanTypeUsage   = "扫描文件的标题语法类型，可以为html，md或multi（multi同时支持所有type）"
	catalogIndentUsage     = "目录的缩进，默认为2空格，输入\\t以替代tab"
	writeBackUsage         = "是否将目录写入文件指定位置"
	tocMarkUsage           = "指定文件中写入目录的位置"
	excludeTitleUsage      = "过滤掉内容等于参数指定值的标题"
	excludeFilterUsage     = "过滤掉内容和参数指定的表达式匹配的标题"
	noExcludeUsage         = "不过滤任何标题"
)

// 字符串参数的默认值
const (
	topTagDefault            = "h2"
	formatterDefault         = "default"
	catalogIdDefault         = "bookmark"
	catalogTitleDefault      = "本文索引"
	catalogOutputTypeDefault = "html"
	catalogScanTypeDefault   = "multi"
	catalogIndentDefault     = "  " // 2 space
	tocMarkDefault           = "[TOC]"
	excludeTitleDefault      = "参考"
	excludeFilterDefault     = ""
)

var usage = fmt.Sprintf(`Usage: %s [option]... <file>

读入file，根据其内容生成目录结构。
未提供file参数时默认读取stdin。

可选参数：
-t string, --top-tag string
	%s (default: "%s")
-f string, --formatter string
	%s (default: "%s")
--catalog-id string
	%s (default: "%s")
--title string
	%s (default: "%s")
-o string, --output string
	%s (default: "%s")
-l string, --title-language string
	%s (default: "%s")
-i string, --indent string
	%s (default: 2空格)
-m string, --toc-mark string
	%s (default: "%s")
--exclude-title=[title1,title2]
	%s (default: "%s")
--exclude-filter=[pattern]
	%s (default: "%s")
-w	%s
--no-exclude	%s
-h, --help	显示本帮助信息并终止程序
`,
	os.Args[0],
	topTagUsage, topTagDefault,
	formatterUsage, formatterDefault,
	catalogIdUsage, catalogIdDefault,
	catalogTitleUsage, catalogTitleDefault,
	catalogOutputTypeUsage, catalogOutputTypeDefault,
	catalogScanTypeUsage, catalogScanTypeDefault,
	catalogIndentUsage,
	tocMarkUsage, tocMarkDefault,
	excludeTitleUsage, excludeTitleDefault,
	excludeFilterUsage, excludeFilterDefault,
	writeBackUsage,
	noExcludeUsage)
