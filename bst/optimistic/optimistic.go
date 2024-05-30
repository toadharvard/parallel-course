package optimistic

import (
	"fmt"
	"sync"

	"golang.org/x/exp/constraints"
)

type node[T constraints.Ordered] struct {
	data  T
	left  *node[T]
	right *node[T]
	mux   sync.Mutex
}

type Tree[T constraints.Ordered] struct {
	root *node[T]
	mux  sync.Mutex
}

func (nd *node[T]) lock() {
	if nd == nil {
		return
	}
	nd.mux.Lock()
}

func (nd *node[T]) unlock() {
	if nd == nil {
		return
	}
	nd.mux.Unlock()
}

func (t *Tree[T]) lock() {
	if t == nil {
		return
	}
	t.mux.Lock()
}

func (t *Tree[T]) unlock() {
	if t == nil {
		return
	}
	t.mux.Unlock()
}

func (t *Tree[T]) validate(data T, nd, parent *node[T]) bool {
	if nd == nil && parent == nil {
		return t.root == nil
	}
	curr := t.root
	var prev *node[T] = nil
	for curr != nil && curr.data != data && curr != nd {
		prev = curr
		if data < curr.data {
			curr = curr.left
		} else {
			curr = curr.right
		}
	}
	return curr == nd && prev == parent
}

func (t *Tree[T]) parentUnlock(parent *node[T]) {
	if parent == nil {
		t.unlock()
	} else {
		parent.unlock()
	}
}

func (t *Tree[T]) findNodeAndParent(data T) (*node[T], *node[T]) {
	for {
		t.lock()
		if t.root == nil {
			return nil, nil
		}
		current := t.root
		var parent *node[T] = nil
		for current != nil && current.data != data {
			grandparent := parent
			parent = current
			if data < current.data {
				current = current.left
			} else {
				current = current.right
			}
			if grandparent == nil {
				t.unlock()
			}
		}
		parent.lock()
		current.lock()
		if t.validate(data, current, parent) {
			return current, parent
		}
		current.unlock()
		parent.unlock()
	}
}

func (t *Tree[T]) Insert(data T) {
	nd, parent := t.findNodeAndParent(data)
	inserted := &node[T]{data: data}
	if t.root == nil {
		defer t.unlock()
		t.root = inserted
		return
	}
	defer t.parentUnlock(parent)
	if nd != nil {
		defer nd.unlock()
		return
	}
	if data < parent.data {
		parent.left = inserted
	} else {
		parent.right = inserted
	}
}

func (t *Tree[T]) Find(data T) bool {
	nd, parent := t.findNodeAndParent(data)
	defer t.parentUnlock(parent)
	if nd != nil {
		defer nd.unlock()
		return true
	}
	return false
}

func (t *Tree[T]) Remove(data T) {
	nd, parent := t.findNodeAndParent(data)
	defer t.parentUnlock(parent)
	if nd == nil {
		return
	}
	defer nd.unlock()

	if nd.left == nil && nd.right == nil {
		if nd == t.root {
			t.root = nil
		} else if nd.data < parent.data {
			parent.left = nil
		} else {
			parent.right = nil
		}
		return
	}

	if nd.left == nil {
		if nd == t.root {
			t.root = nd.right
		} else if nd.data < parent.data {
			parent.left = nd.right
		} else {
			parent.right = nd.right
		}
		return
	}

	if nd.right == nil {
		if nd == t.root {
			t.root = nd.left
		} else if nd.data < parent.data {
			parent.left = nd.left
		} else {
			parent.right = nd.left
		}
		return
	}

	nd.right.lock()
	succParent := nd
	succ := nd.right
	for succ.left != nil {
		succGrandparent := succParent
		succParent = succ
		succ.left.lock()
		succ = succ.left
		if succGrandparent != nil && succGrandparent != nd {
			succGrandparent.unlock()
		}
	}

	defer succ.unlock()
	if succParent != nd {
		defer succParent.unlock()
		succParent.left = succ.right
	} else {
		succParent.right = succ.right
	}
	nd.data = succ.data
}

// Concurrent unsafe methods
func (t *Tree[T]) String() string {
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

func (t *Tree[T]) IsValid() bool {
	return t.root.IsValid()
}

func (nd *node[T]) IsValid() bool {
	if nd == nil {
		return true
	}
	if nd.left != nil && nd.left.data > nd.data {
		return false
	}
	if nd.right != nil && nd.right.data < nd.data {
		return false
	}
	return nd.left.IsValid() && nd.right.IsValid()
}

func New[T constraints.Ordered]() *Tree[T] {
	return &Tree[T]{}
}

func (t *Tree[T]) IsEmpty() bool {
	return t.root == nil
}
