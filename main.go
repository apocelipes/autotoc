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

	catalogIndent := StringFlagWithShortName("indent",
		"i",
		"  ",
		"目录的缩进，默认为2空格(使用prettyprint或output为md时不支持)")

	writeBack := flag.Bool("w", false, "是否将目录写入文件指定位置")
	tocMark := StringFlagWithShortName("toc-mark",
		"m",
		"[TOC]",
		"指定文件中写入目录的位置")

	flag.Parse()
	if len(flag.Args()) == 0 {
		fmt.Fprint(os.Stderr, "错误：需要一个输入文件。\n")
		os.Exit(1)
	}
	// 终端可能无法直接输入tab，所以用\t代替
	if *catalogIndent == "\\t" {
		*catalogIndent = "\t"
	}

	var err error
	var f *os.File
	if !*writeBack {
		f, err = os.Open(flag.Arg(0))
	} else {
		f, err = os.OpenFile(flag.Arg(0), os.O_RDWR, 0644)
	}
	if err != nil {
		panic(err)
	}
	defer f.Close()

	ret := parser.MarkdownParser(f, *topTag, *catalogScanType)
	var data string
	switch *catalogOutputType {
	case "html":
		html := strings.Builder{}
		for _, v := range ret {
			html.WriteString(v.Html())
		}
		data = format.RenderCatalog(*catalogId, *catalogTitle, html.String())

		formatHtmlFunc := format.NewFormatter(*formatter)
		data = formatHtmlFunc(data, *catalogIndent)
	case "md":
		md := strings.Builder{}
		md.WriteString("#### " + *catalogTitle + ":\n")
		for _, v := range ret {
			// each parent has no indent
			md.WriteString(v.Markdown(*catalogIndent, true))
		}

		data = md.String()
	}

	if !*writeBack {
		fmt.Println(data)
		return
	}

	err = WriteBackFile(data, *tocMark, f)
	if err != nil {
		panic(err)
	}
}
