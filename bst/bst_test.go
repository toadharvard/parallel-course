package bst_test

import (
	"bst/coarsegrained"
	"bst/finegrained"
	"bst/optimistic"
	"math/rand"
	"sync"
	"testing"
)

var values []int
var numberSize = 100_000

type BST interface {
	Insert(x int)
	Find(x int) bool
	Remove(x int)
	IsValid() bool
	String() string
	IsEmpty() bool
}

func init() {
	rnd := rand.New(rand.NewSource(0xdeadbeef))
	values = make([]int, numberSize)
	for i := range values {
		values[i] = rnd.Intn(1000)
	}
}

func TestFuzzing(t *testing.T) {
	trees := map[string]BST{
		"CoarseGrained": coarsegrained.New[int](),
		"FineGrained":   finegrained.New[int](),
		"Optimistic":    optimistic.New[int](),
	}
	var wg sync.WaitGroup

	for name, tree := range trees {
		t.Run("Insert"+name, func(t *testing.T) {
			for _, v := range values {
				wg.Add(1)
				go func(x int) {
					defer wg.Done()
					tree.Insert(x)
				}(v)
			}
			wg.Wait()
			if !tree.IsValid() {
				t.Errorf("tree is not valid")
			}
		})

		t.Run("FindAfterInsert"+name, func(t *testing.T) {
			for _, v := range values {
				wg.Add(1)
				go func(x int) {
					defer wg.Done()
					tree.Insert(x)
				}(v)
			}
			wg.Wait()
			for _, v := range values {
				if !tree.Find(v) {
					t.Errorf("tree is not valid")
				}
			}
		})

		t.Run("RemoveAfterInsert"+name, func(t *testing.T) {
			for _, v := range values {
				wg.Add(1)
				go func(x int) {
					defer wg.Done()
					tree.Insert(x)
				}(v)
			}
			wg.Wait()
			for _, v := range values {
				wg.Add(1)
				go func(x int) {
					defer wg.Done()
					tree.Remove(x)
				}(v)
			}
			wg.Wait()
			if !tree.IsEmpty() {
				t.Errorf("tree is not valid")
			}
		})
	}
}
