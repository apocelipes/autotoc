package defaultformatter

import (
	"strings"

	"golang.org/x/net/html"
)

// HtmlElement html节点树
type HtmlElement struct {
	StartToken *html.Token
	EndToken   *html.Token
	SelfClosed bool
	Children   []*HtmlElement
}

// NewHtmlElement 根据startTagToken生成html node
func NewHtmlElement(start *html.Token) *HtmlElement {
	return &HtmlElement{
		StartToken: start,
		EndToken:   nil,
		Children:   make([]*HtmlElement, 0),
	}
}

func (element *HtmlElement) AddChild(child *HtmlElement) {
	if child == nil {
		return
	}

	element.Children = append(element.Children, child)
}

// ToString 将节点树转换成格式化后的字符串
// 只对children的string进行缩进
// 当前节点的缩进应该由上层节点处理
// depth是当前节点的嵌套深度，如果不需要换行则不会增加，从0开始
// subDepth是当前行已经存在的节点的数量，用于控制换行，从1开始
func (element *HtmlElement) ToString(indent string, depth, subDepth int) string {
	// 为root时
	if element.StartToken == nil {
		if len(element.Children) == 0 {
			return ""
		} else {
			ret := strings.Builder{}
			for _, c := range element.Children {
				ret.WriteString(c.ToString(indent, depth, subDepth) + "\n")
			}

			return ret.String()
		}
	}

	// 处理doctype和纯文本，直接返回
	if element.StartToken.Type == html.DoctypeToken {
		return element.StartToken.String()
	}

	if element.StartToken.Type == html.TextToken {
		return element.StartToken.String()
	}

	// 自闭合tag不存在子节点
	if element.SelfClosed {
		return element.StartToken.String()
	}

	// 如果缺失endTag就根据startTag补充一个
	if element.EndToken == nil {
		element.EndToken = &html.Token{
			Type:     html.EndTagToken,
			DataAtom: element.StartToken.DataAtom,
			Data:     element.StartToken.Data,
			Attr:     nil,
		}
	}

	if len(element.Children) == 0 {
		return element.StartToken.String() + element.EndToken.String()
	}

	ret := strings.Builder{}

	if len(element.Children) > 1 || element.NeedNewLine() {
		ret.WriteString(element.StartToken.String() + "\n")
		for _, c := range element.Children {
			ret.WriteString(strings.Repeat(indent, depth+1))
			ret.WriteString(c.ToString(indent, depth+1, 1) + "\n")
		}
		ret.WriteString(strings.Repeat(indent, depth))
		ret.WriteString(element.EndToken.String())

		return ret.String()
	}

	// 一行的node不能超过三个，超过的则从下一行开始
	if subDepth > 3 {
		ret.WriteString("\n" + strings.Repeat(indent, depth))
		ret.WriteString(element.StartToken.String() + "\n")
		ret.WriteString(strings.Repeat(indent, depth+1))
		ret.WriteString(element.Children[0].ToString(indent, depth, 1) + "\n")
		ret.WriteString(strings.Repeat(indent, depth))
		ret.WriteString(element.EndToken.String() + "\n")
		ret.WriteString(strings.Repeat(indent, depth))
	} else {
		ret.WriteString(element.StartToken.String())
		ret.WriteString(element.Children[0].ToString(indent, depth, subDepth+1))
		ret.WriteString(element.EndToken.String())
	}

	return ret.String()
}

// NeedNewLine 如果一个node不管是否有子节点都需要换行，返回true
func (element *HtmlElement) NeedNewLine() bool {
	switch element.StartToken.Data {
	case "body":
		fallthrough
	case "ul":
		fallthrough
	case "ol":
		fallthrough
	case "html":
		return true
	}

	return false
}
