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
	flag.Usage = func() {
		fmt.Fprintln(flag.CommandLine.Output(), usage)
		os.Exit(1)
	}

	topTag := StringFlagWithShortName("top-tag",
		"t",
		topTagDefault,
		topTagUsage)

	formatter := StringFlagWithShortName("formatter",
		"f",
		formatterDefault,
		formatterUsage)

	catalogId := flag.String("catalog-id",
		catalogIdDefault,
		catalogIdUsage)
	catalogTitle := flag.String("title",
		catalogTitleDefault,
		catalogTitleUsage)

	catalogOutputType := StringFlagWithShortName("output",
		"o",
		catalogOutputTypeDefault,
		catalogOutputTypeUsage)

	catalogScanType := StringFlagWithShortName("title-language",
		"l",
		catalogScanTypeDefault,
		catalogScanTypeUsage)

	catalogIndent := StringFlagWithShortName("indent",
		"i",
		catalogIndentDefault,
		catalogIndentUsage)

	writeBack := flag.Bool("w", false, writeBackUsage)
	tocMark := StringFlagWithShortName("toc-mark",
		"m",
		tocMarkDefault,
		tocMarkUsage)

	flag.Parse()
	if len(flag.Args()) == 0 {
		fmt.Fprint(os.Stderr, "错误：需要一个输入文件。\n")
		flag.Usage()
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
