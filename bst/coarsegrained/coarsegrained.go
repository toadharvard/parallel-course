package coarsegrained

import (
	"bst/sequential"
	"sync"

	"golang.org/x/exp/constraints"
)

type Tree[T constraints.Ordered] struct {
	seq *sequential.Tree[T]
	mux sync.Mutex
}

func New[T constraints.Ordered]() *Tree[T] {
	return &Tree[T]{
		seq: sequential.New[T](),
	}
}

func (t *Tree[T]) Insert(data T) {
	t.mux.Lock()
	defer t.mux.Unlock()
	t.seq.Insert(data)
}

func (t *Tree[T]) Remove(data T) {
	t.mux.Lock()
	defer t.mux.Unlock()
	t.seq.Remove(data)
}

func (t *Tree[T]) Find(data T) bool {
	t.mux.Lock()
	defer t.mux.Unlock()
	return t.seq.Find(data)
}

func (t *Tree[T]) IsValid() bool {
	t.mux.Lock()
	defer t.mux.Unlock()
	return t.seq.IsValid()
}

func (t *Tree[T]) String() string {
	t.mux.Lock()
	defer t.mux.Unlock()
	return t.seq.String()
}

func (t *Tree[T]) IsEmpty() bool {
	t.mux.Lock()
	defer t.mux.Unlock()
	return t.seq.IsEmpty()
}
