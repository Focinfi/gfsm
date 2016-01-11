package gfsm

import (
	"testing"
)

const (
	orderPending = "pending"
	orderPaid    = "paid"
	orderShipped = "shipped"
)

var stateMachine = NewStateMachine()

type Order struct {
	name string
	*State
}

var order = &Order{name: "order"}

func TestNewState(t *testing.T) {
	order.State = NewState(orderPending, stateMachine)
	if order.CurrentState() != orderPending {
		t.Errorf("can not new a state, current state is:%s", order.CurrentState())
	}
}

func TestAddEvent(t *testing.T) {
	stateMachine.AddEvent("pay")
	if _, ok := stateMachine.events["pay"]; !ok {
		t.Error("can not add pay")
	}
}

func TestAddTransition(t *testing.T) {
	stateMachine.AddEvent("pay").AddTransition(orderPaid, orderPending)
	event := stateMachine.events["pay"]
	if event.to != orderPaid || event.froms[0] != orderPending {
		t.Errorf("can not add a cerrect transition, current is: to %s, from %v", event.to, event.froms)
	}
}

func TestTrigger(t *testing.T) {
	if err := order.Trigger("pay"); err != nil {
		t.Errorf("can not trigger event pay, err is: %v", err.Error())
	}

	if order.CurrentState() != orderPaid {
		t.Errorf("can not trigger the event pay, current state is: %s", order.CurrentState())
	}
}

func TestTriggerUnknownEvent(t *testing.T) {
	order.currentState = orderPending

	if err := order.Trigger("unknown event"); err == nil {
		t.Error("can not stop the unknown event")
	}

	if order.CurrentState() != orderPending {
		t.Errorf("changed the currentState to: %s", order.CurrentState())
	}
}

func TestTriggerStoppedEvent(t *testing.T) {
	order.currentState = orderShipped

	if err := order.Trigger("pay"); err == nil {
		t.Error("can not stop the unsupported event")
	} else {
		t.Log(err.Error())
	}

	if order.CurrentState() != orderShipped {
		t.Errorf("changed the currentState to: %s", order.CurrentState())
	}
}
