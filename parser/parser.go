package parser

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
)

// MarkdownParser 根据topTag和scanType解析文件中的标题行
func MarkdownParser(file io.Reader, topTag, scanType string) []*TitleNode {
	scanner := bufio.NewScanner(file)
	ret := make([]*TitleNode, 0)

	var mdParser *regexp.Regexp
	switch scanType {
	case "html":
		mdParser = getHtmlParser(topTag)
	case "md":
		mdParser = getRawParser(topTag)
	}
	if mdParser == nil {
		panic("topTag: " + topTag + "not support")
	}

	var parent *TitleNode
	for scanner.Scan() {
		line := scanner.Text()
		if mdParser.MatchString(line) && !strings.HasPrefix(line, "##### 参考") {
			var node *TitleNode
			switch scanType {
			case "html":
				node = parseHtmlTitle(mdParser, line)
			case "md":
				node = parseRawTitle(mdParser, line)
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

// parseHtmlTitle 解析markdown中的html标签
// 将file中的<h1>...</h1>等标题解析成TitleNode
// topTag为作为目录顶层项的html tag名，例如"h2"
// file中的标题行需要为单独的一行，且只能是h1-h5中的某一个tag，
// 例如：<h2 id="title">Title</h2>
// 传入不支持的tag name将会导致panic
func parseHtmlTitle(reg *regexp.Regexp, line string) *TitleNode {
	tagName := reg.FindStringSubmatch(line)[1]
	id := reg.FindStringSubmatch(line)[2]
	content := reg.FindStringSubmatch(line)[3]

	return NewTitleNode(tagName, id, content)
}

// 根据topTag生成相应的正则表达式
func getHtmlParser(topTag string) *regexp.Regexp {
	titleTagParser := regexp.MustCompile(`^h([1-5])$`)
	if !titleTagParser.MatchString(topTag) {
		return nil
	}

	titleSize := titleTagParser.FindStringSubmatch(topTag)[1]
	regFormat := `^<(h[%s-5]) id="(.+)">(.+)</h[%s-5]>$`
	reg := fmt.Sprintf(regFormat, titleSize, titleSize)
	return regexp.MustCompile(reg)
}

// parseRawTitle 解析markdown语法标题为TitleNode
// 标题需要为单独的一行，例如"## Title"
// FIXME 中文字符作为html tag的id时未做浏览器兼容性测试，谨慎使用
func parseRawTitle(reg *regexp.Regexp, line string) *TitleNode {
	tag := reg.FindStringSubmatch(line)[1]
	tagName := fmt.Sprintf("h%d", strings.Count(tag, "#"))
	// markdown语法的标题自动添加与标题内容相同的id作为锚
	content := reg.FindStringSubmatch(line)[2]
	id := content

	return NewTitleNode(tagName, id, content)
}

// 生成markdown语法标题对应的正则
func getRawParser(topTag string) *regexp.Regexp {
	if len(topTag) != 2 {
		return nil
	}

	titleSize, err := strconv.Atoi(string(topTag[1]))
	if err != nil {
		return nil
	}

	reg := fmt.Sprintf(`^(#{%d,5}) (.+)$`, titleSize)
	return regexp.MustCompile(reg)
}
