package parser

import (
	"bufio"
	"fmt"
	"io"
)

// ContentEncoder 对URL进行编码的函数，传入URL，传出编码后的结果
type ContentEncoder func(string) string

// Parser 进行整个markdown的解析
type Parser struct {
	// 解析title
	innerParser TitleParser
	// 开始解析的标题层次
	topTag string
	// 待扫描的标题行的类型
	scanType string
	// 从tocMark后的位置开始解析标题
	tocMark string
	// 根据规则或标题内容过滤title
	filter TitleFilter
	// 对标题的URL进行encoding
	encoder ContentEncoder
}

// TitleParser 解析文件中的标题结构
type TitleParser interface {
	// content如果不是标题就返回nil，否则返回TitleNode
	Parse(content string) *TitleNode
}

// Creator 创建parser的函数类型，接受一个string作为顶层标题结构，被SetParser调用
type Creator func(string) TitleParser

// 存储所有TitleParser的creator
var parserCreators = make(map[string]Creator)

// SetInnerParser 注册新的标题解析器
func SetInnerParser(name string, creator Creator) {
	if name == "" || creator == nil {
		panic("name or creator should not be empty")
	}

	parserCreators[name] = creator
}

// GetInnerParser 获取标题解析器。topTag是可标题的最大级别（例如h1-h5中的h1）
func GetInnerParser(name, topTag string) (TitleParser, error) {
	creator, ok := parserCreators[name]
	if !ok {
		return nil, fmt.Errorf("scanType: %s not support", name)
	}

	return creator(topTag), nil
}

// GetParser 根据name创建markdown文档解析器
func GetParser(options ...Option) *Parser {
	parser := &Parser{}
	for _, o := range options {
		o(parser)
	}
	var err error
	parser.innerParser, err = GetInnerParser(parser.scanType, parser.topTag)
	if err != nil {
		panic(err)
	}

	return parser
}

type Option func(*Parser)

// WithTopTag 设置开始解析的标题层级，值为h1-h5
func WithTopTag(tag string) Option {
	return func(p *Parser) {
		p.topTag = tag
	}
}

// WithScanType 设置需要解析的标题类型（html, markdown）
func WithScanType(t string) Option {
	return func(p *Parser) {
		p.scanType = t
	}
}

// TocMark 设置文档中的toc标记
func WithTOCMark(mark string) Option {
	return func(p *Parser) {
		p.tocMark = mark
	}
}

// WithFilter 根据标题内容确定是否需要被解析
func WithFilter(filter TitleFilter) Option {
	return func(p *Parser) {
		p.filter = filter
	}
}

// WithURLEncoder 设置标题内容的编码器，非ascii字符作为URL时可能需要编码
func WithURLEncoder(encoder ContentEncoder) Option {
	return func(p *Parser) {
		p.encoder = encoder
	}
}

func (p *Parser) encode(url string) string {
	if p.encoder == nil {
		return url
	}

	return p.encoder(url)
}

func (p *Parser) isFiltered(content string) bool {
	if p.filter == nil {
		return false
	}

	return p.filter.FilterTitleContent(content)
}

// Parse 根据topTag和scanType解析文件中的标题行
// 默认从tocMark之后的位置开始解析，如果没有找到tocMark则解析整个markdown文件
func (p *Parser) Parse(file io.Reader) []*TitleNode {
	scanner := bufio.NewScanner(file)
	ret := make([]*TitleNode, 0)

	parents := NewParentStack[*TitleNode]()
	for scanner.Scan() {
		line := scanner.Text()
		// 如果遇到tocMark就重新构建节点树
		if line == p.tocMark {
			ret = make([]*TitleNode, 0)
			parents.Clear()
			continue
		}

		node := p.innerParser.Parse(line)
		if node == nil || p.isFiltered(node.content) {
			continue
		}

		node.id = p.encode(node.id)

		parent, _ := parents.Top()

		// 找到自己的父节点，排除所有同级或下级标题
		for parent != nil && !parent.IsChildNode(node) {
			parents.Pop()
			parent, _ = parents.Top()
		}

		// 因为不知道是否存在子节点，所以都当做拥有子节点并入栈
		parents.Push(node)
		if parent == nil {
			ret = append(ret, node)
			continue
		}

		parent.AddChild(node)
	}

	return ret
}
