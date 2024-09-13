package parser

import (
	"fmt"
	"strings"

	"github.com/apocelipes/autotoc/internal/utils"
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

// WriteHTML 生成当前节点及其子节点的html
func (t *TitleNode) WriteHTML(builder *strings.Builder) {
	if !t.hasChild() {
		fmt.Fprintf(builder, catalogItemTemplate, t.id, t.content)
		return
	}

	fmt.Fprintf(builder, catalogItemWithChildren, t.id, t.content)
	builder.WriteString("<ul>\n")
	for _, child := range t.children {
		child.WriteHTML(builder)
	}
	builder.WriteString("</ul>\n")
	builder.WriteString("</li>\n")
}

// WriteMarkdown 生成当前节点及其子节点的markdown
// depth为0表示当前节点没有父节点，因此不设置缩进，只设置其子节点的缩进
func (t *TitleNode) WriteMarkdown(builder *strings.Builder, indent string, depth int) {
	if depth > 0 {
		utils.RepeatToBuilder(builder, indent, depth)
	}
	fmt.Fprintf(builder, catalogMarkdownTemplate, t.content, t.id)

	if !t.hasChild() {
		return
	}

	for _, child := range t.children {
		child.WriteMarkdown(builder, indent, depth+1)
	}
}
