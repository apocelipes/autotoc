package parser

import (
	"bufio"
	"fmt"
	"io"
)

// TitleParser 解析文件中的标题结构
type TitleParser interface {
	// content如果不是标题就返回nil，否则返回TitleNode
	Parse(content string) *TitleNode
}

// Creator 创建parser的函数类型，接受一个string作为顶层标题结构，被SetParser调用
type Creator func(string) TitleParser

// 存储所有TitleParser的creator
var parserCreators = make(map[string]Creator)

// SetParser 注册新的解析器
func SetParser(name string, creator Creator) {
	if name == "" || creator == nil {
		panic("name or creator should not be empty")
	}

	parserCreators[name] = creator
}

// GetParser 根据name创建解析器
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

	parents := NewParentStack()
	for scanner.Scan() {
		line := scanner.Text()
		// 如果遇到tocMark就重新构建节点树
		if line == option.TocMark {
			ret = make([]*TitleNode, 0)
			parents.Clear()
			continue
		}

		node := parser.Parse(line)
		if node != nil {
			if hasFilter && option.Filter.FilterTitleContent(node.content) {
				continue
			}

			if option.Encoder != nil {
				node.id = option.Encoder(node.id)
			}

			parent, _ := parents.Top().(*TitleNode)

			// 找到自己的父节点，排除所有同级或下级标题
			for parent != nil && !parent.IsChildNode(node) {
				parents.Pop()
				parent, _ = parents.Top().(*TitleNode)
			}

			// 因为不知道是否存在子节点，所以都当做拥有子节点并入栈
			parents.Push(node)
			if parent == nil {
				ret = append(ret, node)
				continue
			}

			parent.AddChild(node)
		}
	}

	return ret
}
