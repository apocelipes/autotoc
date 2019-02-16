package parser

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

// 解析markdown格式的title
type markdownTitleParser struct {
	topTag string
	reg    *regexp.Regexp
}

func init() {
	SetParser("md", newMarkdownTitleParser)
}

func newMarkdownTitleParser(topTag string) TitleParser {
	parser := &markdownTitleParser{topTag: topTag}
	parser.reg = parser.getRegexp()
	if parser.reg == nil {
		return nil
	}

	return parser
}

func (m *markdownTitleParser) Parse(content string) *TitleNode {
	if !m.reg.MatchString(content) {
		return nil
	}

	return m.parseTitle(content)
}

// 生成markdown语法标题对应的正则
func (m *markdownTitleParser) getRegexp() *regexp.Regexp {
	if len(m.topTag) != 2 {
		return nil
	}

	titleSize, err := strconv.Atoi(string(m.topTag[1]))
	if err != nil {
		return nil
	}

	reg := fmt.Sprintf(`^(#{%d,5}) (.+)$`, titleSize)
	return regexp.MustCompile(reg)
}

// parseTitle 解析markdown语法标题为TitleNode
// 标题需要为单独的一行，例如"## Title"
func (m *markdownTitleParser) parseTitle(line string) *TitleNode {
	tag := m.reg.FindStringSubmatch(line)[1]
	tagName := fmt.Sprintf("h%d", strings.Count(tag, "#"))
	// markdown语法的标题自动添加与标题内容相同的id作为锚
	content := m.reg.FindStringSubmatch(line)[2]

	// 如果content含有非字母/数字内容，则删除，多余一个的space字符删减至一个空格并替换为“-”
	id := strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return '-'
		}

		if !unicode.IsLetter(r) && !unicode.IsNumber(r) && r != '-' {
			return -1
		}

		return r
	}, content)
	id = regexp.MustCompile(`-{2,}`).ReplaceAllString(id, "-")

	return NewTitleNode(tagName, id, content)
}
