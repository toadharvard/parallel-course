package sequential

import (
	"fmt"

	"golang.org/x/exp/constraints"
)

type node[T constraints.Ordered] struct {
	left  *node[T]
	data  T
	right *node[T]
}
type Tree[T constraints.Ordered] struct {
	root *node[T]
}

func (tree *Tree[T]) findNodeAndParent(data T) (current *node[T], parent *node[T]) {
	current = tree.root
	for current != nil && current.data != data {
		parent = current
		if data < current.data {
			current = current.left
		} else {
			current = current.right
		}
	}
	return
}

func (tree *Tree[T]) Insert(data T) {
	nd, parent := tree.findNodeAndParent(data)
	inserted := &node[T]{data: data}
	if tree.root == nil {
		tree.root = inserted
		return
	}
	if nd != nil {
		return
	}
	if data < parent.data {
		parent.left = inserted
	} else {
		parent.right = inserted
	}
}

func (tree *Tree[T]) Find(data T) bool {
	curr, _ := tree.findNodeAndParent(data)
	return curr != nil
}

func (tree *Tree[T]) Remove(data T) {
	nd, parent := tree.findNodeAndParent(data)
	if nd == nil {
		return
	}

	if nd.left == nil && nd.right == nil {
		if nd == tree.root {
			tree.root = nil
		} else if nd.data < parent.data {
			parent.left = nil
		} else {
			parent.right = nil
		}
		return
	}

	if nd.left == nil {
		if nd == tree.root {
			tree.root = nd.right
		} else if nd.data < parent.data {
			parent.left = nd.right
		} else {
			parent.right = nd.right
		}
		return
	}

	if nd.right == nil {
		if nd == tree.root {
			tree.root = nd.left
		} else if nd.data < parent.data {
			parent.left = nd.left
		} else {
			parent.right = nd.left
		}
		return
	}

	succParent := nd
	succ := nd.right
	for succ.left != nil {
		succParent = succ
		succ = succ.left
	}
	if succParent != nd {
		succParent.left = succ.right
	} else {
		succParent.right = succ.right
	}
	nd.data = succ.data
}

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
