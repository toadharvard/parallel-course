package finegrained

import (
	"fmt"
	"sync"

	"golang.org/x/exp/constraints"
)

type node[T constraints.Ordered] struct {
	left  *node[T]
	data  T
	right *node[T]
	mux   sync.Mutex
}

func (nd *node[T]) lock() {
	if nd != nil {
		nd.mux.Lock()
	}
}
func (nd *node[T]) unlock() {
	if nd != nil {
		nd.mux.Unlock()
	}
}

type FineGrainedTree[T constraints.Ordered] struct {
	root *node[T]
	mux  sync.Mutex
}

func (t *FineGrainedTree[T]) lock() {
	if t != nil {
		t.mux.Lock()
	}
}

func (t *FineGrainedTree[T]) unlock() {
	if t != nil {
		t.mux.Unlock()
	}
}

func New[T constraints.Ordered]() *FineGrainedTree[T] {
	return &FineGrainedTree[T]{}
}

func (t *FineGrainedTree[T]) Valid() bool {
	t.lock()
	defer t.unlock()
	return t.root.valid()
}

func (t *FineGrainedTree[T]) String() string {
	t.lock()
	defer t.unlock()
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

func (t *FineGrainedTree[T]) findNodeAndParent(data T) (nd *node[T], parent *node[T]) {
	t.lock()
	t.root.lock()
	current := t.root

	for current != nil && current.data != data {
		grandparent := parent
		parent = current
		if data < current.data {
			current = current.left
		} else {
			current = current.right
		}
		current.lock()
		t.unlockParent(grandparent)
	}
	return
}

func (t *FineGrainedTree[T]) unlockParent(parent *node[T]) {
	if parent != nil {
		parent.unlock()
	} else {
		t.unlock()
	}
}

func (t *FineGrainedTree[T]) Search(data T) bool {
	nd, parent := t.findNodeAndParent(data)
	found := nd != nil
	nd.unlock()
	t.unlockParent(parent)
	return found
}

func (t *FineGrainedTree[T]) Insert(data T) {
	inserted := &node[T]{data: data}
	nd, parent := t.findNodeAndParent(data)
	if parent != nil {
		if data < parent.data {
			parent.left = inserted
		} else {
			parent.right = inserted
		}
	} else {
		t.root = inserted
	}
	nd.unlock()
	t.unlockParent(parent)
}

func (t *FineGrainedTree[T]) Delete(data T) {
	nd, parent := t.findNodeAndParent(data)
	if nd == nil {
		t.unlockParent(parent)
		return
	}
	defer nd.unlock()
	defer t.unlockParent(parent)

	// Zero children
	if nd.left == nil && nd.right == nil {
		if nd.data < parent.data {
			parent.left = nil
		} else {
			parent.right = nil
		}
		return
	}

	// Two children
	if nd.left != nil && nd.right != nil {
		nd.left.lock()
		nd.right.lock()

		subparent := nd.left
		current := nd.left.right
		for current != nil {
			subgrandparent := subparent
			subparent = current
			current = current.right
			current.lock()
			t.unlockParent(subgrandparent)
		}
		subparent.right = nd.right
		subparent.unlock()

		if nd.data < parent.data {
			parent.left = nd.left
		} else {
			parent.right = nd.left
		}
		nd.left.unlock()
		nd.right.unlock()
		return
	}

	// One child
	if nd.left != nil {
		if nd.data < parent.data {
			parent.left = nd.left
		} else {
			parent.right = nd.left
		}
	} else {
		if nd.data < parent.data {
			parent.left = nd.right
		} else {
			parent.right = nd.right
		}
	}
}
