package parser

import (
	"fmt"
	"strings"
)

// TitleNode 标题节点，按顶层标签构建的节点树
// 节点树的生长方向是单一的固定的，
// 假设节点按从左向右生长则最新添加的子节点一定在最右或者是最右子节点的后代
type TitleNode struct {
	// html的tag，id，inner text
	// 或markdown格式的标题，id和content相同，tagName与'#'的数量对应
	tagName string
	id      string
	content string

	children []*TitleNode
}

// NewTitleNode 创建新节点，tag是html tag名称
func NewTitleNode(tag, id, content string) *TitleNode {
	return &TitleNode{
		tagName:  tag,
		id:       id,
		content:  content,
		children: make([]*TitleNode, 0),
	}
}

// AddChild 添加子节点
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

// IsChildNode 检查node是否是当前节点的子节点
// 按tagName比较：h1 > h2 > ... > h5
// 因为markdown被顺序解析，所以tag更小意味着它是当前节点的子节点
func (t *TitleNode) IsChildNode(node *TitleNode) bool {
	// 字符串按字典序比较，与比较规则相反
	return t.tagName < node.tagName
}

func (t *TitleNode) hasChild() bool {
	return len(t.children) != 0
}

// HTML 生成当前节点及其子节点的html
func (t *TitleNode) HTML() string {
	if !t.hasChild() {
		return fmt.Sprintf(catalogItemTemplate, t.id, t.content)
	}

	html := strings.Builder{}
	html.WriteString(fmt.Sprintf(catalogItemWithChildren, t.id, t.content))
	html.WriteString("<ul>\n")
	for _, child := range t.children {
		html.WriteString(child.HTML())
	}
	html.WriteString("</ul>\n")
	html.WriteString("</li>\n")

	return html.String()
}

// Markdown 生成当前节点及其子节点的markdown
// depth为0表示当前节点没有父节点，因此不设置缩进，只设置其子节点的缩进
func (t *TitleNode) Markdown(indent string, depth int) string {
	if !t.hasChild() {
		if depth == 0 {
			return fmt.Sprintf(catalogMarkdownTemplate, t.content, t.id)
		}

		return fmt.Sprintf(indent+catalogMarkdownTemplate, t.content, t.id)
	}

	md := strings.Builder{}
	if depth == 0 {
		md.WriteString(fmt.Sprintf(catalogMarkdownTemplate, t.content, t.id))
	} else {
		md.WriteString(fmt.Sprintf(indent+catalogMarkdownTemplate, t.content, t.id))
	}

	for _, child := range t.children {
		md.WriteString(child.Markdown(strings.Repeat(indent, depth+1), depth+1))
	}

	return md.String()
}
