package event_test

import (
	"github.com/stretchr/testify/assert"
	"gomc/src/event"
	"testing"
)

type myTestEvent struct {
	myTestString string
	myTestInt    int
}

func TestBus_RegisterListener(t *testing.T) {
	bus := event.NewBus()

	var (
		receivedTestEvent *myTestEvent
	)
	bus.RegisterListener(func(testEvent *myTestEvent) error {
		receivedTestEvent = testEvent
		return nil
	})
	evt := &myTestEvent{
		myTestString: "Hello, world!",
		myTestInt:    42,
	}
	err := bus.Emit(evt)
	assert.Nil(t, err)
	assert.Equal(t, evt, receivedTestEvent)
}
