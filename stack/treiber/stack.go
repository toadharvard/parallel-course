package treiber

import (
	"errors"
	"sync/atomic"
)

type stack[T any] struct {
	head atomic.Pointer[node[T]]
}

type node[T any] struct {
	value T
	next  atomic.Pointer[node[T]]
}

func (stack *stack[T]) Pop() (nilVar T, Err error) {
	for {
		head := stack.head.Load()
		if head == nil {
			return nilVar, errors.New("nil pointer to stack")
		}
		if stack.head.CompareAndSwap(head, head.next.Load()) {
			return head.value, nil
		}
	}
}

func (stack *stack[T]) Push(val T) {
	newHead := node[T]{value: val}
	for {
		head := stack.head.Load()
		newHead.next.Store(head)
		if stack.head.CompareAndSwap(head, &newHead) {
			return
		}
	}
}

func (stack *stack[T]) Size() int {
	elemCounter := 0
	if stack == nil || stack.head.Load() == nil {
		return 0
	}
	currHead := stack.head.Load()
	for currHead != nil {
		elemCounter++
		currHead = currHead.next.Load()
	}
	return elemCounter
}

func NewStack[T any]() stack[T] {
	return stack[T]{}
}
