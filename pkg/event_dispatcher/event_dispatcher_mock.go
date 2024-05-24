package event_dispatcher

import (
	"context"
	_ "github.com/kr/pretty"
	"github.com/stretchr/testify/mock"
	"sync"
)

type MockDispatcher struct {
	mock.Mock
}

func NewMockDispatcher() *MockDispatcher {
	return &MockDispatcher{}
}

func (m *MockDispatcher) Register(eventName string, handler EventHandlerInterface) error {
	return nil
}
func (m *MockDispatcher) DispatchAndWait(ctx *context.Context, event EventInterface) error {
	m.Called(event)
	return nil
}
func (m *MockDispatcher) Dispatch(ctx *context.Context, event EventInterface, wg *sync.WaitGroup) error {
	//fmt.Printf("%# v\n", pretty.Formatter(event.GetPayload()))
	m.Called(event)
	return nil
}
func (m *MockDispatcher) Remove(eventName string, handler EventHandlerInterface) error {
	return nil
}
func (m *MockDispatcher) Has(eventName string, handler EventHandlerInterface) bool {
	return true
}
func (m *MockDispatcher) Clear() {

}
