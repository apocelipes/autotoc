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

func init() {
	SetParser("html", newHtmlTitleParser)
}

func newHtmlTitleParser(topTag string) TitleParser {
	parser := &htmlTitleParser{topTag: topTag}
	parser.reg = parser.getRegexp()
	if parser.reg == nil {
		return nil
	}

	return parser
}

func (h *htmlTitleParser) Parse(content string) *TitleNode {
	if !h.reg.MatchString(content) {
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

// 根据topTag生成相应的正则表达式
func (h *htmlTitleParser) getRegexp() *regexp.Regexp {
	titleTagParser := regexp.MustCompile(`^h([1-5])$`)
	if !titleTagParser.MatchString(h.topTag) {
		return nil
	}

	titleSize := titleTagParser.FindStringSubmatch(h.topTag)[1]
	regFormat := `^<(h[%s-5]) id="(.+?)">(.+)</h[%s-5]>$`
	reg := fmt.Sprintf(regFormat, titleSize, titleSize)
	return regexp.MustCompile(reg)
}
