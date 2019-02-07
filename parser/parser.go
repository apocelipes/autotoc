package parser

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
)

// MarkdownParser 将file中的<h1>...</h1>等标题解析成TitleNode
// topTag为作为目录顶层项的html tag名，例如"h2"
// file中的标题行需要为单独的一行，且只能是h1-h5中的某一个tag，
// 例如：<h2 id="title">Title</h2>
// 传入不支持的tag name将会导致panic
func MarkdownParser(file io.Reader, topTag string) []*TitleNode {
	scanner := bufio.NewScanner(file)
	ret := make([]*TitleNode, 0)
	parserTitleHtml := getParser(topTag)
	if parserTitleHtml == nil {
		panic("topTag: "+topTag+"not support")
	}
	var parent *TitleNode

	for scanner.Scan() {
		line := scanner.Text()
		if parserTitleHtml.MatchString(line) {
			tagName := parserTitleHtml.FindStringSubmatch(line)[1]
			id := parserTitleHtml.FindStringSubmatch(line)[2]
			content := parserTitleHtml.FindStringSubmatch(line)[3]
			node := NewTitleNode(tagName, id, content)

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

// 根据topTag生成相应的正则表达式
func getParser(topTag string) *regexp.Regexp {
	titleTagParser := regexp.MustCompile(`^h([1-5])$`)
	if !titleTagParser.MatchString(topTag) {
		return nil
	}

	titleSize := titleTagParser.FindStringSubmatch(topTag)[1]
	regFormat := `^<(h[%s-5]) id="(.+)">(.+)</h[%s-5]>$`
	reg := fmt.Sprintf(regFormat, titleSize, titleSize)
	return regexp.MustCompile(reg)
}
