package coarsegrained

import (
	"fmt"
	"sync"

	"golang.org/x/exp/constraints"
)

type node[T constraints.Ordered] struct {
	left  *node[T]
	data  T
	right *node[T]
}

type CoarseGrainedTree[T constraints.Ordered] struct {
	root *node[T]
	mux  sync.Mutex
}

func New[T constraints.Ordered]() *CoarseGrainedTree[T] {
	return &CoarseGrainedTree[T]{}
}

func (t *CoarseGrainedTree[T]) Insert(data T) {
	t.mux.Lock()
	t.root = t.root.insert(data)
	t.mux.Unlock()
}

func (nd *node[T]) insert(data T) *node[T] {
	if nd == nil {
		return &node[T]{
			data: data,
		}
	}
	if data < nd.data {
		nd.left = nd.left.insert(data)
	} else if data > nd.data {
		nd.right = nd.right.insert(data)
	}
	return nd
}

func (t *CoarseGrainedTree[T]) Delete(data T) {
	t.mux.Lock()
	t.root = t.root.delete(data)
	t.mux.Unlock()
}

func (nd *node[T]) delete(data T) *node[T] {
	if nd == nil {
		return nil
	}
	if data < nd.data {
		nd.left.delete(data)
	} else if data > nd.data {
		nd.right.delete(data)
	} else {
		if nd.left == nil || nd.right == nil {
			if nd.left == nil {
				return nd.right
			} else {
				return nd.left
			}
		}

		temp := nd.right.min()
		nd.data = temp.data
		nd.right = nd.right.delete(temp.data)
	}
	return nd
}

func (t *CoarseGrainedTree[T]) Search(data T) bool {
	t.mux.Lock()
	defer t.mux.Unlock()
	return t.root.search(data) != nil
}

func (nd *node[T]) search(data T) *node[T] {
	if nd == nil {
		return nil
	}
	if data < nd.data {
		return nd.left.search(data)
	} else if data > nd.data {
		return nd.right.search(data)
	} else {
		return nd
	}
}

func (nd *node[T]) min() *node[T] {
	current := nd
	for current.left != nil {
		current = current.left
	}
	return current
}

func (t *CoarseGrainedTree[T]) Valid() bool {
	t.mux.Lock()
	defer t.mux.Unlock()
	return t.root.valid()
}

func (nd *node[T]) valid() bool {
	if nd == nil {
		return true
	}
	if nd.left != nil && nd.left.data > nd.data {
		return false
	}
	if nd.right != nil && nd.right.data < nd.data {
		return false
	}
	return nd.left.valid() && nd.right.valid()
}

func (t *CoarseGrainedTree[T]) String() string {
	t.mux.Lock()
	defer t.mux.Unlock()
	return fmt.Sprintf("%v", t.root.inorder())
}

func (nd *node[T]) inorder() []T {
	if nd == nil {
		return nil
	}
	left := nd.left.inorder()
	right := nd.right.inorder()
	return append(append(left, nd.data), right...)
}
