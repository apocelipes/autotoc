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

	// markdown的锚标签为经过处理后的content
	id := m.filterId(content)

	return NewTitleNode(tagName, id, content)
}

// 将id过滤为符合gfm的格式
func (m *markdownTitleParser) filterId(id string) string {
	// gfm下无法显示html tag的效果，并且tag会导致id的值不可预测，所以不支持解析
	if regexp.MustCompile(`</?([a-zA-Z]|[a-zA-Z]+[0-9]+)+ ?.*?>`).MatchString(id) {
		panic("not support html tags in a markdown style title")
	}

	// 如果content含有非字母/数字/下划线内容，则删除，多余一个的space字符删减至一个空格并替换为“-”
	id = strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return '-'
		}

		if !unicode.IsLetter(r) && !unicode.IsNumber(r) && r != '-' && r != '_' {
			return -1
		}

		return r
	}, id)
	id = regexp.MustCompile(`-{2,}`).ReplaceAllString(id, "-")

	// 分别去除每个单词上用于表示斜体字体的下划线，如"__word__"
	words := strings.Split(id, "-")
	for i := range words {
		words[i] = delEnclosedWordUnderline(words[i])
	}

	id = strings.Join(words, "-")
	// 删除开头和结尾配对的下划线
	return delEnclosedWordUnderline(id)
}

// 找到从左起最后一个连续下划线的index
func findLeftUnderline(runes []rune) int {
	length := len(runes)
	for i := 0; i < length; i++ {
		if runes[i] != '_' {
			return i - 1
		}
	}

	return -1
}

// 找到从右起最后一个连续下划线的index
func findRightUnderline(runes []rune) int {
	length := len(runes)
	if length == 0 {
		return -1
	}

	for i := length - 1; i >= 0; i-- {
		if runes[i] != '_' {
			return i + 1
		}
	}

	return -1
}

// 删除单词左右两边的下划线，下划线左右成对删除，多余的保留
func delEnclosedWordUnderline(s string) string {
	// 处理unicode字符
	data := []rune(s)
	// 找到最左和最右的下划线的坐标，如果s只有下划线组成/一边没有下划线，则返回原字符串
	left := findLeftUnderline(data)
	right := findRightUnderline(data)
	if left == -1 || right == -1 || right == len(data) {
		return s
	}

	// 左右下划线连续的长度，用于配对删除
	leftLen := left + 1
	rightLen := len(data) - right

	// 如果左边的连续下划线比右边长，则全部删除右边，删除左边和右边长度一样的部分，其余保留
	// 一样长则左右全部删除
	word := make([]rune, 0, len(data)-2)
	if leftLen < rightLen {
		word = append(word, data[left+1:len(data)-left]...)
	} else if leftLen > rightLen {
		word = append(word, data[left:right]...)
	} else {
		word = append(word, data[left+1:right]...)
	}

	return string(word)
}
