package nonconcurent

import "errors"

type stack[T any] struct {
	head *node[T]
}

type node[T any] struct {
	value T
	next  *node[T]
}

func (stack *stack[T]) Pop() (T, error) {
	if stack.head == nil {
		return *new(T), errors.New("stack is already empty")
	}
	lastValue := stack.head.value
	stack.head = stack.head.next
	return lastValue, nil
}

func (stack *stack[T]) Push(val T) {
	newNode := node[T]{value: val}
	stack.head, newNode.next = &newNode, stack.head
}

func (stack *stack[T]) Size() int {
	elemCounter := 0
	if stack == nil || stack.head == nil {
		return 0
	}
	currHead := stack.head
	for currHead != nil {
		elemCounter++
		currHead = currHead.next
	}
	return elemCounter
}

func NewStack[T any]() stack[T] {
	return stack[T]{}
}
