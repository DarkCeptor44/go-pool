package pool

import (
	"math/rand"
	"sync"
	"time"
)

type Pool struct {
	numWorkers  int
	wg          sync.WaitGroup
	task        Task
	values      []Value
	valuesLock  sync.Mutex
	results     []Result
	resultsLock sync.Mutex
}

type Value struct {
	Val interface{}
}

type Result struct {
	Old Value
	New Value
}

// Pool task has to be this format
type Task func(index int, value Value) Result

func init() {
	rand.NewSource(time.Now().UnixNano())
}

// Returns a new Pool struct
func NewPool(numWorkers int, task Task, values []Value) *Pool {
	return &Pool{
		numWorkers: numWorkers,
		wg:         sync.WaitGroup{},
		task:       task,
		values:     values,
		results:    make([]Result, 0),
	}
}

// Executes the task on the pool with amount of workers specified and returns a list of results
func (p *Pool) Run() []Result {
	for {
		p.valuesLock.Lock()
		if len(p.values) == 0 {
			p.valuesLock.Unlock()
			break
		}

		workerCount := min(p.numWorkers, len(p.values))
		for i := 0; i < workerCount; i++ {
			index := rand.Intn(len(p.values))
			val := p.values[index]
			p.values = removeIndex(p.values, index)

			p.wg.Add(1)
			go func(index int, value Value) {
				defer p.wg.Done()
				res := p.task(index, value)

				p.resultsLock.Lock()
				p.results = append(p.results, res)
				p.resultsLock.Unlock()
			}(i, val)
		}
		p.valuesLock.Unlock()
	}

	p.wg.Wait()

	return p.results
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func removeIndex(arr []Value, index int) []Value {
	result := make([]Value, len(arr)-1)
	copy(result, arr[:index])
	copy(result[index:], arr[index+1:])
	return result
}

// to Value from primitives

// Turns a generic slice into []Value
func FromSlice[K any](slice []K) []Value {
	result := make([]Value, len(slice))
	for i, v := range slice {
		result[i] = Value{Val: v}
	}
	return result
}

// from Value to primitives

// Turns Value into string
func (v *Value) String() string {
	return v.Val.(string)
}

// Turns Value into int
func (v *Value) Int() int {
	return v.Val.(int)
}

// Turns Value into float64
func (v *Value) Float() float64 {
	return v.Val.(float64)
}
