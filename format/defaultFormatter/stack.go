package defaultFormatter

// ParentStack 用于解析html时缓存
type ParentStack struct {
	parents []*HtmlElement
}

// NewParentStack 返回ParentStack，设置第一个元素为nil，用作扫描的初始化状态
func NewParentStack() *ParentStack {
	return &ParentStack{
		parents: make([]*HtmlElement, 1),
	}
}

func (s *ParentStack) Push(e *HtmlElement) {
	s.parents = append(s.parents, e)
}

func (s *ParentStack) Pop() {
	s.parents = s.parents[:len(s.parents)-1]
}

func (s *ParentStack) Top() *HtmlElement {
	return s.parents[len(s.parents)-1]
}
