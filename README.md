# gofun
FP(Functional-Programming) API for Go.

# Usage

```Go
import "fun"

fun.MakeFun(func(handler fun.Handler) {
    input := bufio.NewScanner(os.Stdin)
    for input.Scan() {
      text := input.Text()
      handler.SendNext(text)
    }
	}).Map(func(value interface{}) interface{} {
	  return person{
	    value.(string),
	  }
	}).SubscribeNext(func(next interface{}) {
	  fmt.Println(next.(person))
	})
```
