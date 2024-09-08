package defaultformatter

import (
	"strings"

	"github.com/apocelipes/autotoc/internal/utils"

	"golang.org/x/net/html"
)

// HtmlElement html节点树
type HtmlElement struct {
	StartToken html.Token
	EndToken   html.Token
	Children   []*HtmlElement
	SelfClosed bool
}

// NewHtmlElement 根据startTagToken生成html node
func NewHtmlElement(start html.Token) *HtmlElement {
	return &HtmlElement{
		StartToken: start,
		Children:   make([]*HtmlElement, 0),
	}
}

func (element *HtmlElement) AddChild(child *HtmlElement) {
	if child == nil {
		return
	}

	element.Children = append(element.Children, child)
}

// WriteTo 将节点树转换成格式化后的字符串
// 只对children的string进行缩进
// 当前节点的缩进应该由上层节点处理
// depth是当前节点的嵌套深度，如果不需要换行则不会增加，从0开始
// subDepth是当前行已经存在的节点的数量，用于控制换行，从1开始
func (element *HtmlElement) WriteTo(builder *strings.Builder, indent string, depth, subDepth int) {
	// 为root时
	// we already excluded the error tokens, so html.ErrorToken only contains empty token
	if element.StartToken.Type == html.ErrorToken {
		for _, c := range element.Children {
			c.WriteTo(builder, indent, depth, subDepth)
			builder.WriteByte('\n')
		}
		return
	}

	// 自闭合tag不存在子节点
	if element.SelfClosed || element.StartToken.Type == html.DoctypeToken || element.StartToken.Type == html.TextToken {
		builder.WriteString(element.StartToken.String())
		return
	}

	// 如果缺失endTag就根据startTag补充一个
	if element.EndToken.Type == html.ErrorToken {
		element.EndToken = html.Token{
			Type:     html.EndTagToken,
			DataAtom: element.StartToken.DataAtom,
			Data:     element.StartToken.Data,
			Attr:     nil,
		}
	}

	if len(element.Children) == 0 {
		builder.WriteString(element.StartToken.String())
		builder.WriteString(element.EndToken.String())
		return
	}

	if len(element.Children) > 1 || element.NeedNewLine() {
		builder.WriteString(element.StartToken.String())
		builder.WriteByte('\n')
		for _, c := range element.Children {
			utils.RepeatToBuilder(builder, indent, depth+1)
			c.WriteTo(builder, indent, depth+1, 1)
			builder.WriteByte('\n')
		}
		utils.RepeatToBuilder(builder, indent, depth)
		builder.WriteString(element.EndToken.String())

		return
	}

	// 一行的node不能超过三个，超过的则从下一行开始
	if subDepth > 3 {
		builder.WriteByte('\n')
		utils.RepeatToBuilder(builder, indent, depth)
		builder.WriteString(element.StartToken.String())
		builder.WriteByte('\n')
		utils.RepeatToBuilder(builder, indent, depth+1)
		element.Children[0].WriteTo(builder, indent, depth, 1)
		builder.WriteByte('\n')
		utils.RepeatToBuilder(builder, indent, depth)
		builder.WriteString(element.EndToken.String())
		builder.WriteByte('\n')
		utils.RepeatToBuilder(builder, indent, depth)
	} else {
		builder.WriteString(element.StartToken.String())
		element.Children[0].WriteTo(builder, indent, depth, subDepth+1)
		builder.WriteString(element.EndToken.String())
	}
}

// NeedNewLine 如果一个node不管是否有子节点都需要换行，返回true
func (element *HtmlElement) NeedNewLine() bool {
	switch element.StartToken.Data {
	case "ul", "ol", "body", "html":
		return true
	}

	return false
}
