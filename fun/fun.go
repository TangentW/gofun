package fun

// Event
type Event int

const (
	EventNext Event = iota
	EventError
)

// handler
type Handler struct {
	handler func(event Event, value interface{})
}

// Send Value
func (handler Handler) SendNext(v interface{}) {
	handler.handler(EventNext, v)
}

func (handler Handler) SendError(err error) {
	handler.handler(EventError, err)
}

// Handle
func (handler *Handler) Handle(todo func(event Event, value interface{})) {
	handler.handler = todo
}

func (handler *Handler) HandleEvents(next func(interface{}), err func(error)) {
	handler.Handle(func(event Event, value interface{}) {
		switch event {
		case EventNext:
			if next == nil {
				break
			}
			next(value)
		case EventError:
			if err == nil {
				break
			}
			err(value.(error))
		}
	})
}

func (handler *Handler) HandleNext(next func(interface{})) {
	handler.HandleEvents(next, nil)
}

func (handler *Handler) Output(output chan<- interface{}) {
  handler.HandleNext(func(value interface{}) {
    output <- value
  })
}

func (handler *Handler) HandleError(err func(error)) {
	handler.HandleEvents(nil, err)
}

// Fun
type Fun struct {
	init func(Handler)
	handlers []Handler
}

func New(init func(Handler)) *Fun {
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

func (fun *Fun) Output(output chan<- interface{}) {
  fun.SubscribeNext(func(value interface{}) {
    output <- value
  })
}

func (fun *Fun) SubscribeError(err func(error)) {
  var handler Handler
  handler.HandleError(err)
  fun.Subscribe(handler)
}

// Monad
func Return(value interface{}) *Fun {
  return New(func(handler Handler) { handler.SendNext(value) })
}

func Error(err error) *Fun {
  return New(func(handler Handler) { handler.SendError(err) })
}

func Never() *Fun {
  return New(func(handler Handler) { })
}

func List(elements ...interface{}) *Fun {
  return New(func(handler Handler) {
    for _, element := range elements {
      handler.SendNext(element)
    }
  })
}

func (fun *Fun) Bind(mapper func(interface{}) *Fun) *Fun {
  ret := New(func(handler Handler) {
    fun.Subscribe(Handler{
      func(event Event, value interface{}) {
        switch event {
        case EventNext: 
          mapper(value).Subscribe(handler)
        case EventError:
          handler.SendError(value.(error))
        }
      },
    })
  })
  return ret
}

func (fun *Fun) Map(mapper func(interface{}) interface{}) *Fun {
  return fun.Bind(func(value interface{}) *Fun {
    return Return(mapper(value))
  })
}

func (fun *Fun) Filter(predicate func(interface{}) bool) *Fun {
  return fun.Bind(func(value interface{}) *Fun {
    if predicate(value) {
      return Return(value)
    } else {
      return Never()
    }
  })
}

func (fun *Fun) Do(todo func(interface{})) *Fun {
  return fun.Bind(func(value interface{}) *Fun {
    todo(value)
    return Return(value)
  })
}

