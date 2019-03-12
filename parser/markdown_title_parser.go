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

var (
	// 匹配html tag
	htmlTagMatcher = regexp.MustCompile(`</?([a-zA-Z]|[a-zA-Z]+[0-9]+)+ ?.*?>`)
	// 匹配全部由下划线组成的字符串
	allUnderlineMatcher = regexp.MustCompile(`^_+$`)
	// 匹配非下划线开头和结尾的，中间包含任意字符（包括下划线）的字符串（通常为单词或字母加数字）
	wordMatcher = regexp.MustCompile(`[^_]+(_+?[^_]+)*`)
)

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

	// 判断topTag是否为h1到h5的tag
	titleSize, err := strconv.Atoi(string(m.topTag[1]))
	if err != nil || strings.ToLower(string(m.topTag[0])) != "h" {
		return nil
	}

	reg := fmt.Sprintf(`^(#{%d,6}) (.+)$`, titleSize)
	return regexp.MustCompile(reg)
}

// parseTitle 解析markdown语法标题为TitleNode
// 标题需要为单独的一行，例如"## Title"
func (m *markdownTitleParser) parseTitle(line string) *TitleNode {
	//TODO 可选参考信息作为独立标题结构或被过滤
	if line == "##### 参考" {
		return nil
	}

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
	if htmlTagMatcher.MatchString(id) {
		panic("not support html tags in a markdown style title")
	}

	id = removeNonLetter(id)

	// 分别去除每个单词上用于表示斜体字体的下划线，如"__word__"
	words := strings.Split(id, "-")
	for i := range words {
		words[i] = removeEnclosedWordUnderline(words[i])
	}

	id = strings.ToLower(strings.Join(words, "-"))
	// 删除开头和结尾配对的下划线
	// 如果开头或结尾的非空字符串全部为下划线则不进行第二次删除
	leftNonEmpty := findLeftNonEmptyStr(words)
	rightNonEmpty := findRightNonEmptyStr(words)
	if leftNonEmpty != -1 && rightNonEmpty != -1 {
		if allUnderlineMatcher.MatchString(words[leftNonEmpty]) ||
			allUnderlineMatcher.MatchString(words[rightNonEmpty]) {
			return id
		}
	}

	return removeEnclosedWordUnderline(id)
}

// 如果s含有非字母/数字/连字符/下划线内容，则删除，多余一个的space字符删减至一个空格并替换为“-”
func removeNonLetter(s string) string {
	ret := strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return '-'
		}

		if !unicode.IsLetter(r) && !unicode.IsNumber(r) && r != '-' && r != '_' {
			return -1
		}

		return r
	}, s)

	// space已经被替换为"-"
	return regexp.MustCompile(`-{2,}`).ReplaceAllString(ret, "-")
}

// 找到从左起第一个非空字符串
func findLeftNonEmptyStr(words []string) int {
	length := len(words)
	for i := 0; i < length; i++ {
		if words[i] != "" {
			return i
		}
	}

	return -1
}

// 找到从右起第一个非空字符串
func findRightNonEmptyStr(words []string) int {
	length := len(words)
	if length == 0 {
		return -1
	}

	for i := length - 1; i >= 0; i-- {
		if words[i] != "" {
			return i
		}
	}

	return -1
}

// 查找最左连续下划线的结束位置和最右连续下划线的开始位置
// left返回-1表示不符合匹配规则或没有左侧连续下划线
// right返回-1表示不符合匹配；返回len(s)表示没有右侧连续下划线
func findUnderline(s string) (left, right int) {
	loc := wordMatcher.FindStringIndex(s)
	if loc == nil {
		return -1, -1
	}

	return loc[0] - 1, loc[1]
}

// 删除单词左右两边的下划线，下划线左右成对删除，多余的保留
func removeEnclosedWordUnderline(s string) string {
	// 找到最左和最右的下划线的坐标，如果s只有下划线组成/一边没有下划线，则返回原字符串
	left, right := findUnderline(s)
	if left == -1 || right == -1 || right == len(s) {
		return s
	}

	// 左右下划线连续的长度，用于配对删除
	leftLen := left + 1
	rightLen := len(s) - right

	// 如果左边的连续下划线比右边长，则全部删除右边，删除左边和右边长度一样的部分，其余保留
	// 一样长则左右全部删除
	word := strings.Builder{}
	if leftLen < rightLen {
		word.WriteString(s[left+1 : len(s)-leftLen])
	} else if leftLen > rightLen {
		word.WriteString(s[left+1-(leftLen-rightLen) : right])
	} else {
		word.WriteString(s[left+1 : right])
	}

	return word.String()
}
