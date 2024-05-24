package event_dispatcher

import (
	"context"
	"sync"
	"time"
)

type EventInterface interface {
	GetName() string
	GetDateTime() time.Time
	GetPayload() interface{}
	SetPayload(payload interface{})
}

type EventHandlerInterface interface {
	Handle(ctx *context.Context, eventFired EventInterface, wg *sync.WaitGroup) error
}

type EventDispatcherInterface interface {
	Register(eventName string, handler EventHandlerInterface) error
	Dispatch(ctx *context.Context, event EventInterface, wg *sync.WaitGroup) error
	DispatchAndWait(ctx *context.Context, event EventInterface) error
	Remove(eventName string, handler EventHandlerInterface) error
	Has(eventName string, handler EventHandlerInterface) bool
	Clear()
}
