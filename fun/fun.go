package fun

/*
import (
  "fmt"
)
*/

// Event
type Event int

const (
	Next Event = iota
	Completed
	Error
)

// handler
type Handler struct {
	handler func(event Event, value interface{})
}

// Send Value
func (handler Handler) SendNext(v interface{}) {
	handler.handler(Next, v)
}

func (handler Handler) SendCompleted() {
	handler.handler(Completed, nil)
}

func (handler Handler) SendError(err error) {
	handler.handler(Error, err)
}

// Handle
func (handler *Handler) Handle(todo func(event Event, value interface{})) {
	handler.handler = todo
}

func (handler *Handler) HandleEvents(next func(interface{}), completed func(), err func(error)) {
	handler.Handle(func(event Event, value interface{}) {
		switch event {
		case Next:
			if next == nil {
				break
			}
			next(value)
		case Completed:
			if completed == nil {
				break
			}
			completed()
		case Error:
			if err == nil {
				break
			}
			err(value.(error))
		}
	})
}

func (handler *Handler) HandleNext(next func(interface{})) {
	handler.HandleEvents(next, nil, nil)
}

func (handler *Handler) HandleCompleted(completed func()) {
	handler.HandleEvents(nil, completed, nil)
}

func (handler *Handler) HandleError(err func(error)) {
	handler.HandleEvents(nil, nil, err)
}

// Fun
type Fun struct {
	init func(Handler)
	handlers []Handler
}

func MakeFun(init func(Handler)) *Fun {
  fun := new(Fun)
  fun.init = init
  return fun
}

func (fun *Fun) activateIfNeeds() {
  if fun.init == nil { return }
  fun.init(Handler{
    func(event Event, value interface{}) {
      for _, handler := range fun.handlers {
        handler.handler(event, value)
      }
    },
  })
  fun.init = nil
}

func (fun *Fun) Subscribe(handler Handler) {
  fun.handlers = append(fun.handlers, handler)
  fun.activateIfNeeds()
}

func (fun *Fun) SubscribeNext(next func(interface{})) {
  var handler Handler
  handler.HandleNext(next)
  fun.Subscribe(handler)
}

func (fun *Fun) SubscribeCompleted(completed func()) {
  var handler Handler
  handler.HandleCompleted(completed)
  fun.Subscribe(handler)
}

func (fun *Fun) SubscribeError(err func(error)) {
  var handler Handler
  handler.HandleError(err)
  fun.Subscribe(handler)
}

// Transform
// Monad
func RetFun(value interface{}) *Fun {
  return MakeFun(func(handler Handler) { handler.SendNext(value) })
}

func (fun *Fun) Bind(mapper func(interface{}) *Fun) *Fun {
  ret := MakeFun(func(handler Handler) {
    h := Handler{
      func(event Event, value interface{}) {
        switch event {
        case Next: 
          hi := Handler{
            func(event Event, value interface{}) {
              switch event {
              case Next:
                handler.SendNext(value)
              case Completed:
                handler.SendCompleted()
              case Error:
                handler.SendError(value.(error))
              }
            },
          }
          mapper(value).Subscribe(hi)
        case Completed:
          handler.SendCompleted()
        case Error:
          handler.SendError(value.(error))
        }
      },
    }
    fun.Subscribe(h)
  })
  return ret
}

func (fun *Fun) Map(mapper func(interface{}) interface{}) *Fun {
  return fun.Bind(func(value interface{}) *Fun {
    return RetFun(mapper(value))
  })
}

