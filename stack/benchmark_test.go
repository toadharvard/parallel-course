package treiberstack

import (
	"runtime"
	"stack/backoff"
	nonconcurent "stack/non-concurent"
	"stack/treiber"
	"sync"
	"testing"
)

const numberOfElements = 1_000_000

type Stack[T any] interface {
	Push(T)
	Pop() (T, error)
	Size() int
}

func PushPopNonConcurent(s Stack[int]) {
	for i := 0; i < numberOfElements; i++ {
		s.Push(i)
	}

	for i := 0; i < numberOfElements; i++ {
		s.Pop()
	}
}

func PushPopConcurrent100(stack Stack[int]) {
	goroutineCount := 100
	wg := sync.WaitGroup{}
	wg.Add(goroutineCount)
	for i := 0; i < goroutineCount; i++ {
		go func() {
			for j := 0; j < (numberOfElements / goroutineCount); j++ {
				stack.Push(j)
			}
			wg.Done()
		}()
	}
	wg.Wait()
	wg.Add(goroutineCount)
	for i := 0; i < goroutineCount; i++ {
		go func() {
			for j := 0; j < numberOfElements/goroutineCount; j++ {
				stack.Pop()
			}
			wg.Done()
		}()
	}
	wg.Wait()
}

func PushPopConcurrentAll(stack Stack[int]) {
	wg := sync.WaitGroup{}
	wg.Add(numberOfElements)
	for j := 0; j < numberOfElements; j++ {
		go func(j int) {
			stack.Push(j)
			wg.Done()
		}(j)
	}
	wg.Wait()

	wg.Add(numberOfElements)

	for j := 0; j < numberOfElements; j++ {
		go func() {
			stack.Pop()
			wg.Done()
		}()
	}
	wg.Wait()
}

func PushPopInRow(s Stack[int]) {
	wg := sync.WaitGroup{}
	wg.Add(200_000)
	for j := 0; j < 200_000; j++ {
		go func(j int) {
			for k := 0; k < 9; k++ {
				s.Push(j)
				s.Pop()
			}
			wg.Done()
		}(j)
	}
	wg.Wait()
}

func PushPopConcurentRand(s Stack[int]) {
	wg := sync.WaitGroup{}
	wg.Add(200_000)
	for j := 0; j < 100_000; j++ {
		go func(j int) {
			for k := 0; k < 9; k++ {
				s.Push(j)
			}
			wg.Done()
		}(j)
		go func() {
			for k := 0; k < 9; k++ {
				s.Pop()
			}
			wg.Done()
		}()
	}
	wg.Wait()
}

func BenchmarkNonConcurrent(b *testing.B) {
	b.Run("Non-concurent:", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			simpleStack := nonconcurent.NewStack[int]()
			PushPopNonConcurent(&simpleStack)
		}
	})

	b.Run("Treiber:", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			treiberStack := treiber.NewStack[int]()
			PushPopNonConcurent(&treiberStack)
		}
	})

	b.Run("Backoff:", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			optimizeTreiberStack := backoff.NewStack[int]()
			PushPopNonConcurent(&optimizeTreiberStack)
		}
	})
}

func BenchmarkConcurrent(b *testing.B) {
	runtime.GOMAXPROCS(16)

	b.Run("Treiber 100:", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			treiberStack := treiber.NewStack[int]()
			PushPopConcurrent100(&treiberStack)
		}
	})

	b.Run("Backoff 100:", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			optimizeTreiberStack := backoff.NewStack[int]()
			PushPopConcurrent100(&optimizeTreiberStack)
		}
	})

	b.Run("Treiber all:", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			treiberStack := treiber.NewStack[int]()
			PushPopConcurrent100(&treiberStack)
		}
	})

	b.Run("Backoff all:", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			optimizeTreiberStack := backoff.NewStack[int]()
			PushPopConcurrent100(&optimizeTreiberStack)
		}
	})
}

func BenchmarkConcurrentElimination(b *testing.B) {
	runtime.GOMAXPROCS(16)

	// b.Run("Non-concurent(RACE):", func(b *testing.B) {
	// 	for i := 0; i < b.N; i++ {
	// 		simpleStack := nonconcurent.NewStack[int]()
	// 		PushPopConcurentRand(&simpleStack)
	// 	}
	// })

	b.Run("Treiber(Rand):", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			treiberStack := treiber.NewStack[int]()
			PushPopConcurentRand(&treiberStack)
		}
	})

	b.Run("Backoff(Rand)", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			optimizeTreiberStack := backoff.NewStack[int]()
			PushPopConcurentRand(&optimizeTreiberStack)
		}
	})

	b.Run("Treiber(Row)", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			treiberStack := treiber.NewStack[int]()
			PushPopInRow(&treiberStack)
		}
	})

	b.Run("Backoff(Row)", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			optimizeTreiberStack := backoff.NewStack[int]()
			PushPopInRow(&optimizeTreiberStack)
		}
	})
}