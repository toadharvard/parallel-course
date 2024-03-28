package backoff

import (
	"errors"
	"sync/atomic"
)

type stack[T any] struct {
	head               atomic.Pointer[node[T]]
	eliminationStorage eliminationStorage[T]
}

type node[T any] struct {
	value T
	next  atomic.Pointer[node[T]]
}

func NewStack[T any]() stack[T] {
	return stack[T]{eliminationStorage: newEliminationStorage[T]()}
}

func (s *stack[T]) Push(value T) {
	newHead := &node[T]{value: value}
	for {
		if s.tryPush(newHead) {
			return
		}
		element, err := s.eliminationStorage.visit(&value)
		if err == nil && element == nil {
			return
		}
	}
}

func (s *stack[T]) Pop() (T, error) {
	for {
		node, err := s.tryPop()
		if err != nil {
			return *new(T), err
		}

		if node != nil {
			return node.value, nil
		}

		element, err := s.eliminationStorage.visit(nil)
		if err == nil && element != nil {
			return *element, nil
		}
	}
}

func (s *stack[T]) tryPush(newHead *node[T]) bool {
	currentHead := s.head.Load()
	newHead.next.Store(currentHead)
	return s.head.CompareAndSwap(currentHead, newHead)
}

func (s *stack[T]) tryPop() (*node[T], error) {
	currentHead := s.head.Load()
	if currentHead == nil {
		return nil, errors.New("empty stack")
	}
	newHead := currentHead.next.Load()
	if s.head.CompareAndSwap(currentHead, newHead) {
		return currentHead, nil
	}
	return nil, nil
}

func (s *stack[T]) Size() int {
	size := 0
	head := s.head.Load()

	for head != nil {
		size++
		head = head.next.Load()
	}

	return size
}
