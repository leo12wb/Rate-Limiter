package event_dispatcher

import (
	"context"
	"errors"
	"github.com/rs/zerolog/log"
	"sync"
	"sync/atomic"
)

var ErrHandlerAlreadyRegistered = errors.New("handler already registered")

type EventDispatcher struct {
	handlers map[string][]EventHandlerInterface
}

func NewEventDispatcher() *EventDispatcher {
	return &EventDispatcher{
		handlers: make(map[string][]EventHandlerInterface),
	}
}

func (ev *EventDispatcher) Dispatch(ctx *context.Context, event EventInterface, wg *sync.WaitGroup) error {
	var accErrors error
	m := sync.Mutex{}
	if handlers, ok := ev.handlers[event.GetName()]; ok {
		for _, handler := range handlers {
			wg.Add(1)
			go func(evtHandler EventHandlerInterface) {
				err := evtHandler.Handle(ctx, event, wg)
				if err != nil {
					log.Error().Err(err).Msg("")
					m.TryLock()
					accErrors = errors.Join(accErrors)
					m.Unlock()
				}
			}(handler)
		}
	}
	return nil
}

func (ev *EventDispatcher) DispatchAndWait(ctx *context.Context, event EventInterface) error {
	var accErrors error
	wg := sync.WaitGroup{}
	var ops atomic.Uint64
	finalRes := make(chan error)
	defer close(finalRes)
	if handlers, ok := ev.handlers[event.GetName()]; ok {
		handlersTotal := len(handlers)
		wg.Add(handlersTotal)
		ch := make(chan error)
		go func() {
			for err := range ch {
				if err != nil {
					log.Error().Err(err).Msg("")
					accErrors = errors.Join(accErrors, err)
				}
			}
			finalRes <- accErrors
		}()
		for _, handler := range handlers {
			go func(evtHandler EventHandlerInterface) {
				ch <- evtHandler.Handle(ctx, event, &wg)
				ops.Add(1)
				if uint64(handlersTotal) == ops.Load() {
					close(ch)
				}
			}(handler)
		}
		wg.Wait()
	}
	accErrors = <-finalRes
	return accErrors
}

func (ed *EventDispatcher) Register(eventName string, handler EventHandlerInterface) error {
	if _, ok := ed.handlers[eventName]; ok {
		for _, h := range ed.handlers[eventName] {
			if h == handler {
				return ErrHandlerAlreadyRegistered
			}
		}
	}
	ed.handlers[eventName] = append(ed.handlers[eventName], handler)
	return nil
}

func (ed *EventDispatcher) Has(eventName string, handler EventHandlerInterface) bool {
	if _, ok := ed.handlers[eventName]; ok {
		for _, h := range ed.handlers[eventName] {
			if h == handler {
				return true
			}
		}
	}
	return false
}

func (ed *EventDispatcher) Remove(eventName string, handler EventHandlerInterface) error {
	if _, ok := ed.handlers[eventName]; ok {
		for i, h := range ed.handlers[eventName] {
			if h == handler {
				ed.handlers[eventName] = append(ed.handlers[eventName][:i], ed.handlers[eventName][i+1:]...)
				return nil
			}
		}
	}
	return nil
}

func (ed *EventDispatcher) Clear() {
	ed.handlers = make(map[string][]EventHandlerInterface)
}
