### gfms

`gfsm` is a finite state machine in Golang.

#### Install

`go get github.com/Focinfi/gfsm`

#### Example

If you are designing a Order struct

```go
  import "github.com/Focinfi/gfsm"

  // design all order's states
  const (
    orderPending = "pending"
    orderPaid    = "paid"
    orderShipped = "shipped"
  )

  // create a new StateMachine
  var OrderStateMachine = gfsm.NewStateMachine(orderPending, orderPaid, orderShipped)

  // add a pay event for OrderStateMachine
  // can only add state among in orderPending, orderPaid, orderShipped, or will return error
  OrderStateMachine.AddEvent("pay").AddTransition(orderPaid, orderPending)
  OrderStateMachine.AddEvent("ship").AddTransition(orderShipped, orderPaid)

  // compose gfsm.State pointer into Order struct
  type Order struct {
    Name string
    *gfsm.State
  }

  // create new object with OrderStateMachine
  var order = &Order{Name: "order", State: gfsm.NewState(orderPending, OrderStateMachine)}

  // get the current state
  order.CurrentState() // "pending"

  // trigger the pay event
  if err := order.Trigger("pay"); err != nil {
    // do something when this order can not pay.
  } else {
    order.CurrentState() // "paid"
    // do something when this order can pay
  }
``` 
