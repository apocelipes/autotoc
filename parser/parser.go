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
func ParseMarkdown(file io.Reader, topTag, scanType string) []*TitleNode {
	scanner := bufio.NewScanner(file)
	ret := make([]*TitleNode, 0)

	parser := GetParser(scanType, topTag)
	if parser == nil {
		panic(fmt.Sprintf("topTag: %s or scanType: %s not support", topTag, scanType))
	}

	var parent *TitleNode
	for scanner.Scan() {
		line := scanner.Text()
		node := parser.Parse(line)
		if node != nil {
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
