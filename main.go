package main

import (
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"

	"github.com/apocelipes/autotoc/format"
	"github.com/apocelipes/autotoc/parser"
)

func checker(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func main() {
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

	var f *os.File
	if len(flag.Args()) == 0 {
		// 未提供文件名参数时判断是否处于pipe中，是则stdin为输入文件，不可能为terminal
		if IsStdinTerminal() {
			_, _ = fmt.Fprintln(os.Stderr, "错误：需要一个输入文件")
			flag.Usage()
		}

		if *writeBack {
			log.Fatalln("-w不能在输入为stdin时使用")
		}

		f = os.Stdin
	} else {
		openFlag := os.O_RDONLY
		if *writeBack {
			openFlag = os.O_RDWR
		}
		var err error
		f, err = os.OpenFile(flag.Arg(0), openFlag, 0)
		checker(err)
		defer f.Close()
	}

	// 终端可能无法直接输入tab，所以用\t代替
	if *catalogIndent == "\\t" {
		*catalogIndent = "\t"
	}

	options := []parser.Option{
		parser.WithTopTag(*topTag),
		parser.WithScanType(*catalogScanType),
		parser.WithTOCMark(*tocMark),
	}

	if !*noExclude {
		titleFilter := &parser.DefaultFilter{}
		titleFilter.SetExcludeTitles(*excludeTitle)
		checker(titleFilter.SetExcludeRegExp(*excludeFilter))
		options = append(options, parser.WithFilter(titleFilter))
	}

	if !*noEncode {
		options = append(options, parser.WithURLEncoder(url.PathEscape))
	}

	mdParser := parser.GetParser(options...)
	ret := mdParser.Parse(f)
	if len(ret) == 0 {
		log.Fatalln("未找到任何标题")
	}

	var data string
	switch *catalogOutputType {
	case "html":
		html := strings.Builder{}
		for _, v := range ret {
			html.WriteString(v.HTML())
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
		//default:
		//	log.Fatalf("unknow format: %v\n", catalogOutputType)
	}

	if *writeBack || *fullOutput {
		checker(WriteCatalog(f, data, *tocMark, *fullOutput))
		return
	}

	fmt.Println(data)
}
