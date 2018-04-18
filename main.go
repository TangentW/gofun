package main

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
  return fmt.Sprintf("Hello, my name is %s\n", p.name)
}

func main() {
	fmt.Println("vim-go")
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

}
