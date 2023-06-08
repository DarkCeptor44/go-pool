package pool

import (
	"math/rand"
	"testing"
	"time"
)

func TestFactorial(t *testing.T) {
	arr := getRandomSlice(10, 20, 50)
	p := NewPool(10, func(index int, value Value) Result {
		n := factorial(value.Int())
		return Result{
			Old: value,
			New: Value{Val: n},
		}
	}, FromSlice[int](arr))

	start := time.Now()
	results := p.Run()
	elapsed := time.Since(start)

	for _, r := range results {
		got := r.New.Int()
		want := factorial(r.Old.Int())

		if want == got {
			t.Logf("%d -> %d\n", r.Old.Int(), got)
		} else {
			t.Errorf("Got %d, wanted %d", got, want)
		}
	}

	t.Logf("Done in %s\n", elapsed)
}

func factorial(n int) int {
	if n == 0 {
		return 1
	}
	return n * factorial(n-1)
}

func init() {
	rand.NewSource(time.Now().UnixNano())
}

// int slice of length o from m to n
func getRandomSlice(m, n, o int) []int {
	result := make([]int, o)
	for i := 0; i < o; i++ {
		result[i] = rand.Intn(n-m+1) + m
	}
	return result
}
