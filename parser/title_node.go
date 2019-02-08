package parser

import (
	"fmt"
	"strings"
)

type TitleNode struct {
	// html的tag，id，inner text
	tagName string
	id      string
	content string

	children []*TitleNode
}

func NewTitleNode(tag, id, content string) *TitleNode {
	return &TitleNode{
		tagName:  tag,
		id:       id,
		content:  content,
		children: make([]*TitleNode, 0),
	}
}

// 添加子节点
// 如果不是直接后代则作为最后一个直接子节点的孩子被添加
func (t *TitleNode) AddChild(child *TitleNode) {
	if !t.hasChild() {
		t.children = append(t.children, child)
		return
	}

	// 因为markdown是被顺序解析的，所以非直接子节点一定是最后一个直接子节点的孩子
	last := t.children[len(t.children)-1]
	if last.IsChildNode(child) {
		last.AddChild(child)
		return
	}

	t.children = append(t.children, child)
}

// 检查node是否是当前节点的子节点
// 按tagName比较：h1 > h2 > ... > h5
// 因为markdown被顺序解析，所以tag更小意味着它是当前节点的子节点
func (t *TitleNode) IsChildNode(node *TitleNode) bool {
	// 字符串按字典序比较，与比较规则相反
	return t.tagName < node.tagName
}

func (t *TitleNode) hasChild() bool {
	return len(t.children) != 0
}

// 生成当前节点及其子节点的html
func (t *TitleNode) Html() string {
	if !t.hasChild() {
		return fmt.Sprintf(catalogItemTemplate, t.id, t.content)
	}

	html := strings.Builder{}
	html.WriteString(fmt.Sprintf(catalogItemWithChildren, t.id, t.content))
	html.WriteString("<ul>\n")
	for _, child := range t.children {
		html.WriteString(child.Html())
	}
	html.WriteString("</ul>\n")
	html.WriteString("</li>\n")

	return html.String()
}

func (t *TitleNode) Markdown(indent string) string {
	if !t.hasChild() {
		return fmt.Sprintf(indent+catalogMarkdownTemplate, t.content, t.id)
	}

	md := strings.Builder{}
	md.WriteString(fmt.Sprintf(indent+catalogMarkdownTemplate, t.content, t.id))
	for _, child := range t.children {
		md.WriteString(child.Markdown("  "+indent))
	}

	return md.String()
}
