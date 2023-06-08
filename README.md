# gopool

Go equivalent-ish of Python's `multiprocessing.Pool`.

## Installation

```bash
go get github.com/DarkCeptor44/go-pool
```

## Example

```go
import (
    "fmt"

    "github.com/DarkCeptor44/go-pool"
)

func multiply(i int, value pool.Value) pool.Result {
 n := value.Int() * 2

 return pool.Result{
  Old: value,
  New: pool.Value{Val: n},
 }
}

func main(){
 arr := []int{1, 2, 3, 4, 5}

 // makes pool with 10 workers, the task function and the slice of ints
 p := pool.NewPool(10, multiply, pool.FromSlice[int](arr))
 results := p.Run()

 for _, result := range results {
    fmt.Printf("Original: %d, Multiplied: %d\n", result.Old.Int(), result.New.Int())
 }
}
```
