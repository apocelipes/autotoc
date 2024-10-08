package defaultformatter

import (
	"fmt"
	"io"
	"strings"

	"github.com/apocelipes/autotoc/internal/stack"

	"golang.org/x/net/html"
)

// parseHtml 将html解析为HtmlElement
// ret是函数的结果集
// stack记录解析时父节点的变化
// 正常解析完成返回io.EOF
func parseHtml(t *html.Tokenizer, ret *HtmlElement, stack *stack.NodeStack[*HtmlElement]) error {
	for {
		t.Next()
		tk := t.Token()

		parent, _ := stack.Top()
		switch tk.Type {
		case html.StartTagToken:
			el := NewHtmlElement(tk)
			if parent != nil {
				parent.AddChild(el)
			} else {
				ret.AddChild(el)
			}
			stack.Push(el)

			if err := parseHtml(t, ret, stack); err != nil {
				return err
			}
		case html.EndTagToken:
			if tk.Data == parent.StartToken.Data {
				// 找到父节点的endTag后返回
				parent.EndToken = tk
				stack.Pop()
				return nil
			}
			return fmt.Errorf("mismatched end tag: %v ends with %v", parent.StartToken.String(), tk.String())
		case html.TextToken:
			stripped := strings.Trim(tk.Data, "\n")
			if len(stripped) == 0 {
				continue
			}
			fallthrough
		case html.CommentToken, html.DoctypeToken, html.SelfClosingTagToken:
			el := NewHtmlElement(tk)
			el.SelfClosed = tk.Type == html.SelfClosingTagToken
			if parent != nil {
				parent.AddChild(el)
			} else {
				ret.AddChild(el)
			}
		case html.ErrorToken:
			return t.Err()
		}
	}
}

// FormatHtml 格式化html
func FormatHtml(data, indent string) (string, error) {
	t := html.NewTokenizer(strings.NewReader(data))
	ret := NewHtmlElement(html.Token{})
	nodeStack := stack.NewNodeStack[*HtmlElement]()
	if err := parseHtml(t, ret, nodeStack); err != nil && err != io.EOF {
		return "", err
	}

	buf := strings.Builder{}
	ret.WriteTo(&buf, indent, 0, 1)
	return buf.String(), nil
}
