package mpdclient

import (
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
	"time"
)

type TestListenerImpl struct {
	someField string
	onEvent   func()
}

func (l TestListenerImpl) OnEvent() {
	l.onEvent()
}

func TestAddOnConnectListener(t *testing.T) {
	t.Run("single listener", func(t *testing.T) {
		listenerImpl := TestListenerImpl{
			someField: "startValue",
		}
		c := make(chan interface{}, 1)
		newValue := "new value"
		listenerImpl.onEvent = func() {
			listenerImpl.someField = newValue
			c <- struct{}{}
		}
		observer := MpdClientImpl{
			host:      testHost,
			port:      testPort,
			password:  testPassword,
			listeners: mpdListeners{},
		}
		observer.AddOnConnectListener(listenerImpl)
		observer.fireOnConnectEvent()
		select {
		case <-c:
		case <-time.After(time.Second):
			t.Fatal("onEvent() was not called in the goroutine")
		}
		assert.Equal(t, listenerImpl.someField, newValue)
	})
	t.Run("multiply listener", func(t *testing.T) {
		t.Log("preparing data")
		wg := sync.WaitGroup{}
		c := make(chan interface{}, 2)
		newValue := "new value"
		t.Log("creating 1st TestListenerImpl instance")
		listenerImpl1 := TestListenerImpl{
			someField: "startValue",
		}
		wg.Add(1)
		listenerImpl1.onEvent = func() {
			listenerImpl1.someField = newValue
			wg.Done()
		}
		t.Log("creating 2nd TestListenerImpl instance")
		listenerImpl2 := TestListenerImpl{
			someField: "startValue",
		}
		wg.Add(1)
		listenerImpl2.onEvent = func() {
			listenerImpl2.someField = newValue
			wg.Done()
		}
		t.Log("creating observer instance")
		observer := MpdClientImpl{
			host:      testHost,
			port:      testPort,
			password:  testPassword,
			listeners: mpdListeners{},
		}
		t.Log("registering listeners")
		observer.AddOnConnectListener(&listenerImpl1)
		observer.AddOnConnectListener(&listenerImpl2)

		t.Log("starting test (fire onConnectEvent)")
		observer.fireOnConnectEvent()

		t.Log("checking results")
		go func() {
			wg.Wait()
			c <- struct{}{}
		}()
		select {
		case <-c:
		case <-time.After(time.Second):
			t.Fatal("onEvent() was not called in the goroutine")
		}
		assert.Equal(t, listenerImpl1.someField, newValue)
		assert.Equal(t, listenerImpl2.someField, newValue)
	})
}

func TestRemoveOnConnectListener(t *testing.T) {
	t.Run("tes", func(t *testing.T) {
		c := make(chan interface{}, 1)
		startValue := "start value"
		newValue := "new value"
		listenerImpl := TestListenerImpl{
			someField: startValue,
		}
		listenerImpl.onEvent = func() {
			listenerImpl.someField = newValue
			c <- struct{}{}
		}
		observer := MpdClientImpl{
			host:      testHost,
			port:      testPort,
			password:  testPassword,
			listeners: mpdListeners{},
		}
		observer.AddOnConnectListener(&listenerImpl)
		observer.RemoveOnConnectListener(&listenerImpl)
		observer.fireOnConnectEvent()
		select {
		case <-c:
			t.Fatal("onEvent() was called in the goroutine")
		case <-time.After(time.Second):
		}
		assert.Equal(t, listenerImpl.someField, startValue)
	})
}

func TestAddOnDisonnectListener(t *testing.T) {
	t.Run("single listener", func(t *testing.T) {
		listenerImpl := TestListenerImpl{
			someField: "startValue",
		}
		c := make(chan interface{}, 1)
		newValue := "new value"
		listenerImpl.onEvent = func() {
			listenerImpl.someField = newValue
			c <- struct{}{}
		}
		observer := MpdClientImpl{
			host:      testHost,
			port:      testPort,
			password:  testPassword,
			listeners: mpdListeners{},
		}
		observer.AddOnDisconnectListener(listenerImpl)
		observer.fireOnDisconnectEvent()
		select {
		case <-c:
		case <-time.After(time.Second):
			t.Fatal("onEvent() was not called in the goroutine")
		}
		assert.Equal(t, listenerImpl.someField, newValue)
	})
	t.Run("multiply listener", func(t *testing.T) {
		t.Log("preparing data")
		wg := sync.WaitGroup{}
		c := make(chan interface{}, 2)
		newValue := "new value"
		t.Log("creating 1st TestListenerImpl instance")
		listenerImpl1 := TestListenerImpl{
			someField: "startValue",
		}
		wg.Add(1)
		listenerImpl1.onEvent = func() {
			listenerImpl1.someField = newValue
			wg.Done()
		}
		t.Log("creating 2nd TestListenerImpl instance")
		listenerImpl2 := TestListenerImpl{
			someField: "startValue",
		}
		wg.Add(1)
		listenerImpl2.onEvent = func() {
			listenerImpl2.someField = newValue
			wg.Done()
		}
		t.Log("creating observer instance")
		observer := MpdClientImpl{
			host:      testHost,
			port:      testPort,
			password:  testPassword,
			listeners: mpdListeners{},
		}
		t.Log("registering listeners")
		observer.AddOnDisconnectListener(&listenerImpl1)
		observer.AddOnDisconnectListener(&listenerImpl2)

		t.Log("starting test (fire onDisconnectEvent)")
		observer.fireOnDisconnectEvent()

		t.Log("checking results")
		go func() {
			wg.Wait()
			c <- struct{}{}
		}()
		select {
		case <-c:
		case <-time.After(time.Second):
			t.Fatal("onEvent() was not called in the goroutine")
		}
		assert.Equal(t, listenerImpl1.someField, newValue)
		assert.Equal(t, listenerImpl2.someField, newValue)
	})
}

func TestRemoveOnDisonnectListener(t *testing.T) {
	t.Run("tes", func(t *testing.T) {
		c := make(chan interface{}, 1)
		startValue := "start value"
		newValue := "new value"
		listenerImpl := TestListenerImpl{
			someField: startValue,
		}
		listenerImpl.onEvent = func() {
			listenerImpl.someField = newValue
			c <- struct{}{}
		}
		observer := MpdClientImpl{
			host:      testHost,
			port:      testPort,
			password:  testPassword,
			listeners: mpdListeners{},
		}
		observer.AddOnDisconnectListener(&listenerImpl)
		observer.RemoveOnDisconnectListener(&listenerImpl)
		observer.fireOnConnectEvent()
		select {
		case <-c:
			t.Fatal("onEvent() was called in the goroutine")
		case <-time.After(time.Second):
		}
		assert.Equal(t, listenerImpl.someField, startValue)
	})
}
