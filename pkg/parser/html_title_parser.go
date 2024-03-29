package parser

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// 解析html格式的title
type htmlTitleParser struct {
	topTag string
	reg    *regexp.Regexp
}

// HTMLTitleParserName name of HTMLTitleParser for the parser factory
const HTMLTitleParserName = "html"

func init() {
	SetInnerParser(HTMLTitleParserName, newHTMLTitleParser)
}

func newHTMLTitleParser(topTag string) TitleParser {
	parser := &htmlTitleParser{topTag: topTag}
	if parser.reg = parser.getRegexp(); parser.reg == nil {
		return nil
	}

	return parser
}

func (h *htmlTitleParser) Parse(content string) *TitleNode {
	if !h.reg.MatchString(content) {
		return nil
	}

	openTag := h.reg.FindStringSubmatch(content)[1]
	closeTag := h.reg.FindStringSubmatch(content)[4]
	if openTag != closeTag {
		return nil
	}

	return h.parseTitle(content)
}

// parseTitle 解析markdown中的html标签
// 将file中的<h1>...</h1>等标题解析成TitleNode
// topTag为作为目录顶层项的html tag名，例如"h2"
// file中的标题行需要为单独的一行，且只能是h1-h5中的某一个tag，
// 例如：<h2 id="title">Title</h2>
// 传入不支持的tag name将会导致panic
func (h *htmlTitleParser) parseTitle(line string) *TitleNode {
	tagName := h.reg.FindStringSubmatch(line)[1]
	id := h.reg.FindStringSubmatch(line)[2]
	// 处理html嵌套tag中的文本内容
	doc, _ := goquery.NewDocumentFromReader(strings.NewReader(line))
	content := doc.Text()

	return NewTitleNode(tagName, id, content)
}

// getRegexp 根据topTag生成相应的正则表达式
func (h *htmlTitleParser) getRegexp() *regexp.Regexp {
	titleTagParser := regexp.MustCompile(`^h([1-6])$`)
	if !titleTagParser.MatchString(h.topTag) {
		return nil
	}

	titleSize := titleTagParser.FindStringSubmatch(h.topTag)[1]
	// golang的regexp不支持反向引用，html tag的开闭标签名必须一样
	regFormat := `^<(h[%s-6])[^>]+?id="(.+?)"[^>]*?>(.+)</(h[%s-6])>$`
	reg := fmt.Sprintf(regFormat, titleSize, titleSize)
	return regexp.MustCompile(reg)
}
