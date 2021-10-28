package model

type Transition struct {
	ID   uint64
	Name string
	From string
	To   string
}

type EventType uint8

type EventStatus uint8

const (
	Created EventType = iota
	Updated
	Removed

	Deferred EventStatus = iota
	Processed
)

type TransitionEvent struct {
	ID     uint64
	Type   EventType
	Status EventStatus
	Entity *Transition
}
