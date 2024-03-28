package backoff

import "math/rand/v2"

type eliminationStorage[T any] struct {
	exchangers []exchanger[T]
	capacity   int
	waitSteps  int
}

func newEliminationStorage[T any]() eliminationStorage[T] {
	result := eliminationStorage[T]{capacity: 10, waitSteps: 1000}
	result.exchangers = make([]exchanger[T], result.capacity)
	for i := 0; i < result.capacity; i++ {
		result.exchangers[i] = newExchanger[T]()
	}
	return result
}

func (es *eliminationStorage[T]) visit(value *T) (*T, error) {
	randRange := func(min, max int) int {
		return rand.IntN(max-min) + min
	}

	index := randRange(0, es.capacity)
	return es.exchangers[index].exchange(value, es.waitSteps)
}
