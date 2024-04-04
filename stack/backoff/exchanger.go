package backoff

import (
	"errors"
	"sync/atomic"
)

type exchangerState int

const (
	empty exchangerState = iota
	waiting
	busy
)

type exchanger[T any] struct {
	item atomic.Value
}

type item[T any] struct {
	value *T
	state exchangerState
}

func newExchanger[T any]() exchanger[T] {
	result := exchanger[T]{}
	result.item.Store(item[T]{nil, empty})
	return result
}

func (ex *exchanger[T]) exchange(value *T, waitSteps int) (*T, error) {
	emptyCase := func(passedSteps int) (*T, error) {
		for j := passedSteps; j < waitSteps; j++ {
			exchangerItem := ex.item.Load().(item[T])
			if exchangerItem.state != busy {
				continue
			}
			newItem := item[T]{nil, empty}
			ex.item.Store(newItem)
			return exchangerItem.value, nil
		}
		return new(T), errors.New("timeout")
	}

	for i := 0; i < waitSteps; i++ {
		exchangerItem := ex.item.Load().(item[T])

		switch exchangerItem.state {
		case empty:
			oldItem := item[T]{nil, empty}
			newItem := item[T]{value, waiting}
			if ex.item.CompareAndSwap(oldItem, newItem) {
				return emptyCase(i)
			}
		case waiting:
			oldItem := item[T]{exchangerItem.value, waiting}
			newItem := item[T]{value, busy}
			if ex.item.CompareAndSwap(oldItem, newItem) {
				return exchangerItem.value, nil
			}
		}
	}
	return new(T), errors.New("timeout")
}
