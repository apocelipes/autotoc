package parser

// ParentStack 用于解析html时缓存
type ParentStack[T any] struct {
	parents []T
}

const defaultCacheSize = 10

// NewParentStack 返回ParentStack，设置第一个元素为nil，用作扫描的初始化状态
func NewParentStack[T any]() *ParentStack[T] {
	return &ParentStack[T]{
		parents: make([]T, 0, defaultCacheSize),
	}
}

func (s *ParentStack[T]) Push(e T) {
	s.parents = append(s.parents, e)
}

func (s *ParentStack[T]) Pop() {
	if len(s.parents) == 0 {
		return
	}

	s.parents = s.parents[:len(s.parents)-1]
}

func (s *ParentStack[T]) Top() (T, bool) {
	if len(s.parents) == 0 {
		var t T
		return t, false
	}

	return s.parents[len(s.parents)-1], true
}

func (s *ParentStack[T]) Clear() {
	s.parents = make([]T, 0, defaultCacheSize)
}
