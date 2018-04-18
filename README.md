# gofun
FP(Functional-Programming) API for Go.

# Usage

```Go
import (
  "fmt"
  "bufio"
  "os"
  "./fun"
)

type person struct {
  name string
}

func (p person) String() string {
  return fmt.Sprintf("Hello, my name is %s", p.name)
}

func main() {
  events := make(chan interface{})

  go obtainInput(events)

  for {
    select {
    case p := <-events:
      fmt.Println(p.(person))
    }
  }
}

func obtainInput(output chan<- interface{}) {
  fun.New(func(handler fun.Handler) {
    input := bufio.NewScanner(os.Stdin)
    for input.Scan() {
      handler.SendNext(input.Text())
    }
  }).Map(func(value interface{}) interface{} {
    return person{
      name: value.(string),
    }
  }).Output(output)
}
```
