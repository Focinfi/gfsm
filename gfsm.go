package gfsm

import (
	"fmt"
	"github.com/Focinfi/gset"
)

type State struct {
	currentState string
	sm           *StateMachine `sql:"-"`
}

func (s State) CurrentState() string {
	return string(s.currentState)
}

type Event struct {
	name  string
	sm    *StateMachine
	froms []string
	to    string
}

type StateMachine struct {
	events map[string]*Event
	states *gset.Set
}

func NewState(currentState string, sm *StateMachine) *State {
	return &State{currentState, sm}
}

func NewStateMachine(states ...string) *StateMachine {
	stateE := make([]gset.Elementer, len(states))
	for i, state := range states {
		stateE[i] = gset.T(state)
	}
	stateSet := gset.NewSet(stateE...)
	return &StateMachine{make(map[string]*Event), stateSet}
}

func (sm *StateMachine) Event(name string) *Event {
	event, ok := sm.events[name]
	if !ok {
		event := &Event{name: name, sm: sm}
		sm.events[name] = event
	}
	return event
}

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
