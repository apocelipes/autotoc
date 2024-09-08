package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/url"
	"os"
	"strings"

	"github.com/apocelipes/autotoc/internal/utils"
	"github.com/apocelipes/autotoc/pkg/format"
	"github.com/apocelipes/autotoc/pkg/parser"
)

func checkError(err error) {
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func main() {
	topTag := utils.StringFlagWithShortName("top-tag",
		"t",
		topTagDefault,
		topTagUsage)

	formatter := utils.StringFlagWithShortName("formatter",
		"f",
		formatterDefault,
		formatterUsage)

	catalogID := flag.String("catalog-id",
		catalogIDDefault,
		catalogIDUsage)
	catalogTitle := flag.String("title",
		catalogTitleDefault,
		catalogTitleUsage)

	catalogOutputType := utils.StringFlagWithShortName("output",
		"o",
		catalogOutputTypeDefault,
		catalogOutputTypeUsage)

	catalogScanType := utils.StringFlagWithShortName("title-language",
		"l",
		catalogScanTypeDefault,
		catalogScanTypeUsage)

	catalogIndent := utils.StringFlagWithShortName("indent",
		"i",
		catalogIndentDefault,
		catalogIndentUsage)

	writeBack := flag.Bool("w", false, writeBackUsage)
	tocMark := utils.StringFlagWithShortName("toc-mark",
		"m",
		tocMarkDefault,
		tocMarkUsage)

	excludeTitle := flag.String("exclude-title",
		excludeTitleDefault,
		excludeTitleUsage)
	excludeFilter := flag.String("exclude-filter",
		"",
		excludeFilterUsage)
	noExclude := flag.Bool("no-exclude", false, noExcludeUsage)

	fullOutput := flag.Bool("full", false, fullOutputUsage)

	noEncode := flag.Bool("no-encode", false, noEncodeUsage)

	flag.Parse()

	if *writeBack && *fullOutput {
		checkError(errors.New("--full 不能和选项 -w 一起使用"))
	}

	var f *os.File
	if len(flag.Args()) == 0 {
		// 未提供文件名参数时判断是否处于pipe中，是则stdin为输入文件，不可能为terminal
		if utils.IsStdinTerminal() {
			_, _ = os.Stderr.WriteString("错误：需要一个输入文件")
			flag.Usage()
		}

		if *writeBack {
			checkError(errors.New("-w 不能在输入为 stdin 时使用"))
		}

		f = os.Stdin
	} else {
		openFlag := os.O_RDONLY
		if *writeBack {
			openFlag = os.O_RDWR
		}
		var err error
		f, err = os.OpenFile(flag.Arg(0), openFlag, 0)
		checkError(err)
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
		checkError(titleFilter.SetExcludeRegExp(*excludeFilter))
		options = append(options, parser.WithFilter(titleFilter))
	}

	if !*noEncode {
		options = append(options, parser.WithURLEncoder(url.PathEscape))
	}

	fileContent, err := io.ReadAll(f)
	checkError(err)
	mdParser := parser.GetParser(options...)
	titleNodes := mdParser.Parse(bytes.NewReader(fileContent))
	if len(titleNodes) == 0 {
		// not an error
		_, _ = os.Stdout.WriteString("未找到任何标题")
		return
	}

	var catalogContent string
	switch *catalogOutputType {
	case "html":
		catalogContent = renderHTMLTitles(*catalogID, *catalogTitle, *catalogIndent, *formatter, titleNodes)
	case "md":
		catalogContent = renderMarkdownTitles(*catalogTitle, *catalogIndent, titleNodes)
	default:
		checkError(fmt.Errorf("不支持的格式化类型: %v", *catalogOutputType))
	}

	if *writeBack {
		checkError(utils.WriteBackFile([]byte(catalogContent), []byte(*tocMark), fileContent, f.Name()))
		return
	} else if *fullOutput {
		checkError(utils.WriteStdout([]byte(catalogContent), []byte(*tocMark), fileContent))
		return
	}

	fmt.Println(catalogContent)
}

func renderHTMLTitles(catalogID, catalogTitle, catalogIndent, formatter string, titles []*parser.TitleNode) string {
	html := strings.Builder{}
	fmt.Fprintf(&html, `<blockquote id="%s"><h4>%s</h4><ul>`, catalogID, catalogTitle)
	for _, v := range titles {
		v.WriteHTML(&html)
	}
	fmt.Fprint(&html, "</ul></blockquote>")

	formatHTMLFunc := format.NewFormatter(formatter)
	if formatHTMLFunc == nil {
		panic(fmt.Errorf("unsupported HTML formatter: %v", formatter))
	}
	return formatHTMLFunc(html.String(), catalogIndent)
}

func renderMarkdownTitles(catalogTitle, catalogIndent string, titles []*parser.TitleNode) string {
	md := strings.Builder{}
	md.WriteString("#### " + catalogTitle + "\n\n")
	for _, v := range titles {
		// each parent has no indent
		v.WriteMarkdown(&md, catalogIndent, 0)
	}

	return md.String()
}
