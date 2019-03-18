package parser

import (
	"bufio"
	"fmt"
	"io"
)

// 解析文件中的标题结构
type TitleParser interface {
	// content如果不是标题就返回nil，否则返回TitleNode
	Parse(content string) *TitleNode
}

// 创建parser的函数类型，接受一个string作为顶层标题结构
type ParserCreator func(string) TitleParser

// 存储所有TitleParser的creator
var parserCreators = make(map[string]ParserCreator)

// 注册新的解析器
func SetParser(name string, creator ParserCreator) {
	if name == "" || creator == nil {
		panic("name or creator should not be empty")
	}

	parserCreators[name] = creator
}

// 根据name创建解析器
func GetParser(name, topTag string) TitleParser {
	creator, ok := parserCreators[name]
	if !ok {
		return nil
	}

	return creator(topTag)
}

// ParseMarkdown 根据topTag和scanType解析文件中的标题行
// 默认从tocMark之后的位置开始解析，如果没有找到tocMark则解析整个markdown文件
func ParseMarkdown(file io.Reader, option *ParseOption) []*TitleNode {
	scanner := bufio.NewScanner(file)
	ret := make([]*TitleNode, 0)
	hasFilter := option.Filter != nil

	parser := GetParser(option.ScanType, option.TopTag)
	if parser == nil {
		panic(fmt.Sprintf("topTag: %s or scanType: %s not support",
			option.TopTag, option.ScanType))
	}

	var parent *TitleNode
	for scanner.Scan() {
		line := scanner.Text()
		// 如果遇到tocMark就重新构建节点树
		if line == option.TocMark {
			ret = make([]*TitleNode, 0)
			parent = nil
			continue
		}

		node := parser.Parse(line)
		if node != nil {
			if hasFilter && option.Filter.FilterTitleContent(node.content) {
				continue
			}

			if parent == nil || !parent.IsChildNode(node) {
				parent = node
				ret = append(ret, node)
				continue
			}

			parent.AddChild(node)
		}
	}

	return ret
}
