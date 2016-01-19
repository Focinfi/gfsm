package gfsm

import (
	"fmt"
	"github.com/Focinfi/gset"
)

// State has pointer of StateMachine tagged with `sql:"-"`
// which means ingore this column when using github.com/jinzhu/gorm
type State struct {
	currentState string
	sm           *StateMachine `sql:"-"`
}

// CurrentState return this state's state
func (s State) CurrentState() string {
	return string(s.currentState)
}

// Event has the pointer of the StateMachine to which this event bolonging
// one event triggerred from some state to one state
type Event struct {
	name  string
	sm    *StateMachine
	froms []string
	to    string
}

// StateMachine has a set of states and events
type StateMachine struct {
	events map[string]*Event
	states *gset.Set
}

// NewState return a pointer of new state
func NewState(currentState string, sm *StateMachine) *State {
	return &State{currentState, sm}
}

// NewStateMachine return a pointer of new StateMachine
func NewStateMachine(states ...string) *StateMachine {
	stateE := make([]gset.Elementer, len(states))
	for i, state := range states {
		stateE[i] = gset.T(state)
	}
	stateSet := gset.NewSet(stateE...)
	return &StateMachine{make(map[string]*Event), stateSet}
}

// Event find or create a event
func (sm *StateMachine) Event(name string) *Event {
	event, ok := sm.events[name]
	if !ok {
		event := &Event{name: name, sm: sm}
		sm.events[name] = event
	}
	return event
}

// Transition add transition for this event, error will not be nil
// if any state in to or froms does not exit
func (e *Event) Transition(to string, froms ...string) error {
	states := append(froms[:], to)

	for _, state := range states {
		if !e.sm.states.Has(gset.T(state)) {
			return fmt.Errorf("has not state: %s", state)
		}
	}

	e.to = to
	e.froms = froms
	return nil
}

// Trigger trigger event, error will not be nil if the event does not exist or
// current state is not in this event's froms
// s.State will be change to the event's to state, when event triggers sucessfully,
// and error is nil
func (s *State) Trigger(eventName string) error {
	event, ok := s.sm.events[eventName]
	if !ok {
		return fmt.Errorf("No such %s event!", eventName)
	}

	for _, stateFrom := range event.froms {
		if s.currentState == stateFrom {
			s.currentState = event.to
			return nil
		}
	}

	return fmt.Errorf("Can not trigger event [%s], from state [%s]", event.name, s.currentState)
}
