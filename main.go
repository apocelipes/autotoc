package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/apocelipes/markdown-catalog-generator/format"
	"github.com/apocelipes/markdown-catalog-generator/parser"
)

func main() {
	topTag := StringFlagWithShortName("top-tag",
		"t",
		"h2",
		"设置作为目录顶层项的tag")

	formatter := StringFlagWithShortName("formatter",
		"f",
		"prettyprint",
		"选择格式化html代码的方式，目前只支持GoHTML和prettyprint(output为markdown时不支持)")

	catalogId := flag.String("catalog-id",
		"bookmark",
		"目录的html id(output为markdown时不支持)")
	catalogTitle := flag.String("title",
		"本文索引",
		"目录的标题")

	catalogOutputType := StringFlagWithShortName("output",
		"o",
		"html",
		"输出的目录格式，可以为html或md(markdown)")

	catalogScanType := StringFlagWithShortName("title-language",
		"l",
		"html",
		"扫描文件的标题语法类型，可以为html或md")

	flag.Parse()
	if len(flag.Args()) == 0 {
		fmt.Fprint(os.Stderr, "错误：需要一个输入文件。\n")
		os.Exit(1)
	}

	f, err := os.Open(flag.Args()[0])
	if err != nil {
		panic(err)
	}
	defer f.Close()

	ret := parser.MarkdownParser(f, *topTag, *catalogScanType)
	switch *catalogOutputType {
	case "html":
		html := strings.Builder{}
		for _, v := range ret {
			html.WriteString(v.Html())
		}
		data := format.RenderCatalog(*catalogId, *catalogTitle, html.String())

		formatHtmlFunc := format.NewFormatter(*formatter)
		fmt.Println(formatHtmlFunc(data))
	case "md":
		md := strings.Builder{}
		md.WriteString("#### "+*catalogTitle+":\n")
		for _, v := range ret {
			// each parent has no indent
			md.WriteString(v.Markdown(""))
		}

		fmt.Println(md.String())
	}
}
