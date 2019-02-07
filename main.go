package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/yosssi/gohtml"

	"github.com/apocelipes/markdown-catalog-generator/parser"
)

func main() {
	topTag := ""
	topTagUsage := "设置作为目录顶层项的tag"
	flag.StringVar(&topTag, "top-tag", "h2", topTagUsage)
	flag.StringVar(&topTag, "t", "h2", topTagUsage)
	flag.Parse()
	if len(flag.Args()) == 0 {
		fmt.Fprint(os.Stderr, "错误：需要一个输入文件。")
		os.Exit(1)
	}

	f, err := os.Open(flag.Args()[0])
	if err != nil {
		panic(err)
	}
	defer f.Close()

	ret := parser.MarkdownParser(f, topTag)
	html := strings.Builder{}
	for _, v := range ret {
		html.WriteString(v.Html())
	}

	fmt.Println(gohtml.Format(html.String()))
}
