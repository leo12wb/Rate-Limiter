package event_dispatcher

import "time"

type DummyEvent struct {
	Name    string
	Payload interface{}
}

func NewDummyEvent() *DummyEvent {
	return &DummyEvent{
		Name: "DummyEvent",
	}
}

func (e *DummyEvent) NewDummyEvent() *DummyEvent {
	return &DummyEvent{
		Name: "DummyEvent",
	}
}

func (e *DummyEvent) GetName() string {
	return e.Name
}

func (e *DummyEvent) GetPayload() interface{} {
	return e.Payload
}

func (e *DummyEvent) SetPayload(payload interface{}) {
	e.Payload = payload
}

func (e *DummyEvent) GetDateTime() time.Time {
	return time.Now()
}
