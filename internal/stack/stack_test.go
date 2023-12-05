package stack

import (
	"slices"
	"testing"
)

func TestNewNodeStack(t *testing.T) {
	intStack := NewNodeStack[int]()
	if cap(intStack.parents) != defaultCacheSize {
		t.Errorf("wrong cap: want %v, got %v", defaultCacheSize, cap(intStack.parents))
	}
	if len(intStack.parents) != 0 {
		t.Errorf("wrong length: want 0, got %v", len(intStack.parents))
	}

	emptyStructStack := NewNodeStack[struct{}]()
	if cap(emptyStructStack.parents) != defaultCacheSize {
		t.Errorf("wrong cap: want %v, got %v", defaultCacheSize, cap(emptyStructStack.parents))
	}
	if len(emptyStructStack.parents) != 0 {
		t.Errorf("wrong length: want 0, got %v", len(emptyStructStack.parents))
	}

	nestStack := NewNodeStack[NodeStack[string]]()
	if cap(nestStack.parents) != defaultCacheSize {
		t.Errorf("wrong cap: want %v, got %v", defaultCacheSize, cap(nestStack.parents))
	}
	if len(nestStack.parents) != 0 {
		t.Errorf("wrong length: want 0, got %v", len(nestStack.parents))
	}
}

func TestStackPushPop(t *testing.T) {
	intStack := NewNodeStack[int]()
	intStack.Push(1)
	if len(intStack.parents) != 1 {
		t.Errorf("wrong length: want 1, got %v", len(intStack.parents))
	}
	if v, ok := intStack.Pop(); !ok || v != 1 {
		t.Errorf("Nodestack.Pop error: want(1, true), got(%v, %v)", v, ok)
	}
	if v, ok := intStack.Pop(); ok || v == 1 {
		t.Errorf("Nodestack.Pop error: want(0, false), got(%v, %v)", v, ok)
	}

	target := []int{10, 9, 8, 7, 6, 5, 4, 3, 2, 1, 0}
	stack := NewNodeStack[int]()
	for _, v := range target {
		stack.Push(v)
	}
	if len(stack.parents) != len(target) {
		t.Errorf("wrong length: want %v, got %v", len(target), len(stack.parents))
	}
	result := make([]int, 0, len(target))
	for {
		v, ok := stack.Pop()
		if !ok {
			break
		}
		result = append(result, v)
	}
	slices.Reverse(target)
	if !slices.Equal(target, result) {
		t.Errorf("Push/Pop failed: want %v, got %v", target, result)
	}
}

func TestStackTop(t *testing.T) {
	intStack := NewNodeStack[int]()
	intStack.Push(1)
	if v, ok := intStack.Top(); !ok || v != 1 {
		t.Errorf("Nodestack.Top error: want(1, true), got(%v, %v)", v, ok)
	}
	intStack.Pop()
	if v, ok := intStack.Top(); ok || v == 1 {
		t.Errorf("Nodestack.Top error: want(0, false), got(%v, %v)", v, ok)
	}
}

func TestStackClear(t *testing.T) {
	stack := NewNodeStack[int]()

	stack.Clear()
	if cap(stack.parents) != defaultCacheSize {
		t.Errorf("wrong cap: want %v, got %v", defaultCacheSize, cap(stack.parents))
	}
	if len(stack.parents) != 0 {
		t.Errorf("wrong length: want 0, got %v", len(stack.parents))
	}

	for i := 0; i < defaultCacheSize+1; i++ {
		stack.Push(i)
	}
	if cap(stack.parents) == defaultCacheSize {
		t.Errorf("wrong cap: got %v", cap(stack.parents))
	}
	if len(stack.parents) != defaultCacheSize+1 {
		t.Errorf("wrong length: want %v, got %v", defaultCacheSize+1, len(stack.parents))
	}
	stack.Clear()
	if cap(stack.parents) != defaultCacheSize {
		t.Errorf("wrong cap: want %v, got %v", defaultCacheSize, cap(stack.parents))
	}
	if len(stack.parents) != 0 {
		t.Errorf("wrong length: want 0, got %v", len(stack.parents))
	}
}
