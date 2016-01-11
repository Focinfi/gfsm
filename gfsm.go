package gfsm

import "fmt"

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
}

func NewState(currentState string, sm *StateMachine) *State {
	return &State{currentState, sm}
}

func NewStateMachine() *StateMachine {
	return &StateMachine{make(map[string]*Event)}
}

func (sm *StateMachine) AddEvent(name string) *Event {
	event, ok := sm.events[name]
	if !ok {
		event := &Event{name: name}
		sm.events[name] = event
	}
	return event
}

func (e *Event) AddTransition(to string, froms ...string) {
	e.to = to
	e.froms = froms
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
