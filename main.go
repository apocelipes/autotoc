package main

import (
	"flag"
	"fmt"
	"net/url"
	"os"
	"strings"

	"golang.org/x/crypto/ssh/terminal"

	"github.com/apocelipes/autotoc/format"
	"github.com/apocelipes/autotoc/parser"
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

	catalogID := flag.String("catalog-id",
		catalogIDDefault,
		catalogIDUsage)
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

	excludeTitle := flag.String("exclude-title",
		excludeTitleDefault,
		excludeTitleUsage)
	excludeFilter := flag.String("exclude-filter",
		excludeFilterDefault,
		excludeFilterUsage)
	noExclude := flag.Bool("no-exclude", false, noExcludeUsage)

	fullOutput := flag.Bool("full", false, fullOutputUsage)

	noEncode := flag.Bool("no-encode", false, noEncodeUsage)

	flag.Parse()

	var err error
	var f *os.File
	if len(flag.Args()) == 0 {
		// 未提供文件名参数时判断是否处于pipe中，是则stdin为输入文件
		if terminal.IsTerminal(int(os.Stdin.Fd())) {
			fmt.Fprint(os.Stderr, "错误：需要一个输入文件。\n")
			flag.Usage()
		}

		if *writeBack {
			fmt.Fprintln(os.Stderr, "-w不能在输入为stdin时使用")
			os.Exit(1)
		}

		f = os.Stdin
	} else {
		if !*writeBack {
			f, err = os.Open(flag.Arg(0))
		} else {
			f, err = os.OpenFile(flag.Arg(0), os.O_RDWR, 0644)
		}
		if err != nil {
			panic(err)
		}
		defer f.Close()
	}

	// 终端可能无法直接输入tab，所以用\t代替
	if *catalogIndent == "\\t" {
		*catalogIndent = "\t"
	}

	option := parser.ParseOption{
		TopTag:   *topTag,
		ScanType: *catalogScanType,
		TocMark:  *tocMark,
	}

	if !*noExclude {
		titleFilter := &parser.DefaultFilter{}
		titleFilter.SetExcludeTitles(*excludeTitle)
		if *excludeFilter != "" {
			if err := titleFilter.SetExcludeRegExp(*excludeFilter); err != nil {
				fmt.Fprintln(os.Stderr, "parse exclude-filter error: "+err.Error())
				os.Exit(1)
			}
		}
		option.Filter = titleFilter
	}

	if !*noEncode {
		option.Encoder = url.PathEscape
	}

	ret := parser.ParseMarkdown(f, &option)
	if len(ret) == 0 {
		fmt.Fprintln(os.Stderr, "未找到任何标题。")
		os.Exit(1)
	}

	var data string
	switch *catalogOutputType {
	case "html":
		html := strings.Builder{}
		for _, v := range ret {
			html.WriteString(v.Html())
		}
		data = format.RenderCatalog(*catalogID, *catalogTitle, html.String())

		formatHTMLFunc := format.NewFormatter(*formatter)
		data = formatHTMLFunc(data, *catalogIndent)
	case "md":
		md := strings.Builder{}
		md.WriteString("#### " + *catalogTitle + "\n\n")
		for _, v := range ret {
			// each parent has no indent
			md.WriteString(v.Markdown(*catalogIndent, 0))
		}

		data = md.String()
	}

	if *writeBack || *fullOutput {
		err = WriteCatalog(f, data, *tocMark, *fullOutput)
		if err != nil {
			panic(err)
		}

		return
	}

	fmt.Println(data)
}
