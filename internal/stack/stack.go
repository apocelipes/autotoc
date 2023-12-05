package stack

// NodeStack 用于解析html时缓存
type NodeStack[T any] struct {
	parents []T
}

const defaultCacheSize = 10

// NewNodeStack 返回ParentStack，设置第一个元素为nil，用作扫描的初始化状态
func NewNodeStack[T any]() *NodeStack[T] {
	return &NodeStack[T]{
		parents: make([]T, 0, defaultCacheSize),
	}
}

func (s *NodeStack[T]) Push(e T) {
	s.parents = append(s.parents, e)
}

func (s *NodeStack[T]) Pop() (T, bool) {
	if len(s.parents) == 0 {
		var value T
		return value, false
	}

	ret := s.parents[len(s.parents)-1]
	s.parents = s.parents[:len(s.parents)-1]
	return ret, true
}

func (s *NodeStack[T]) Top() (T, bool) {
	if len(s.parents) == 0 {
		var t T
		return t, false
	}

	return s.parents[len(s.parents)-1], true
}

func (s *NodeStack[T]) Clear() {
	s.parents = make([]T, 0, defaultCacheSize)
}
