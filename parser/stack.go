package parser

// ParentStack 用于解析html时缓存
type ParentStack struct {
	parents []interface{}
}

// NewParentStack 返回ParentStack，设置第一个元素为nil，用作扫描的初始化状态
func NewParentStack() *ParentStack {
	return &ParentStack{
		parents: make([]interface{}, 0),
	}
}

func (s *ParentStack) Push(e interface{}) {
	s.parents = append(s.parents, e)
}

func (s *ParentStack) Pop() {
	if len(s.parents) == 0 {
		return
	}

	s.parents = s.parents[:len(s.parents)-1]
}

func (s *ParentStack) Top() interface{} {
	if len(s.parents) == 0 {
		return nil
	}

	return s.parents[len(s.parents)-1]
}

func (s *ParentStack) Clear() {
	if len(s.parents) == 0 {
		return
	}

	s.parents = make([]interface{}, 0)
}
