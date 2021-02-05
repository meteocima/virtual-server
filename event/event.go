package event

import (
	"sync"
)

// Source is an interface that represents
// any object which emits one or more events.
type Source interface{}

// Event is a struct that contains infgormations about
// a single event emitted. It contains two fields:
//
// `Source`, which contains the Object which calls `Invoke` method,
// and `Payload` whic contains informations specifics to the event type.
type Event struct {
	Source  Source
	Payload interface{}
}

// Emitter is an object which can emits
// multiple events of the same kind,
//
// Event can be listened to by multiple listeners.
//
// Both events emission and listening can happen in different goroutines.
// The implementation synchronizes all accesses to internal
// `listeners` field using a channel of `listenerAction`
// structs.
type Emitter struct {
	closed     bool
	closedLock *sync.Mutex
	source     Source
	listeners  map[*Listener]struct{}
}

// NewEmitter return a new instance
// of Emitter, linked `source` argument.
func NewEmitter(source Source) *Emitter {
	return &Emitter{
		closedLock: &sync.Mutex{},
		source:     source,
		listeners:  map[*Listener]struct{}{},
	}
}

// InitSource initializes a list of fields
// of a `Source` object.
func InitSource(source Source, emitters ...**Emitter) {
	for _, emitter := range emitters {
		*emitter = NewEmitter(source)
	}
}

// CloseEmitters closes all emitters of a source.
func CloseEmitters(emitters ...**Emitter) {
	for _, emitter := range emitters {
		(*emitter).Close()
		//*emitter = nil
	}
}

// Listener is a single listener
// which is listening on events of an `Emitter`.
type Listener struct {
	c          chan *Event
	closed     bool
	closedLock *sync.Mutex
	e          *Emitter
}

// Stop listening new events invoked
// on a single listener.
func (l *Listener) Stop() {
	actionsOnListeners <- removeListenerAction(l, l.e)
}

func (l *Listener) killChannel() {
	l.closedLock.Lock()
	if !l.closed {
		l.closed = true
		close(l.c)
	}
	l.closedLock.Unlock()
}

// Invoke causes all listeners of this
// `Emitter` to receive an `Event` instance with
// `Emitter` source as source, and `payload` argument
// as `payload`
func (e *Emitter) Invoke(payload interface{}) {
	event := Event{
		Source:  e.source,
		Payload: payload,
	}
	action := emitListenerAction(&event, e)
	actionsOnListeners <- action
	<-action.countResp
}

// Handler ...
type Handler func(ev *Event)

// Count returns total number of
// listeners currently listening
// on this emitter.
func (e *Emitter) Count() int {
	action := countListenersAction(e)
	actionsOnListeners <- action
	return <-action.countResp
}

// Listen add a listener to the emitter that
// executes a function for each event emitted.
func (e *Emitter) Listen(fn Handler) *Listener {
	lst := e.AddListener()
	go func() {
		for event := range lst.c {
			fn(event)
		}
	}()
	return lst
}

// Close removes all listener of
// the `Emitter` calling `Clear` method, and then
// closes the `actionsOnListeners` channel, making the internal `Emitter`
// goroutine to terminate politely
func (e *Emitter) Close() {
	e.closedLock.Lock()
	e.closed = true
	e.closedLock.Unlock()
	actionsOnListeners <- closeEmitterListenersAction(e)
}

// IsClosed ...
func (e *Emitter) IsClosed() bool {
	e.closedLock.Lock()
	defer e.closedLock.Unlock()
	return e.closed
}

// AwaitOne registers a new `Listener` on an
// Emitter, then waits for an event to occurs
// on it, and finally unregisters the listener
// instance.
func (e *Emitter) AwaitOne() *Event {
	//fmt.Printf("AwaitOne %v\n", e)
	if e.IsClosed() {
		return nil
	}
	lst := e.AddListener()
	event := <-lst.c
	lst.Stop()
	return event
}

// AwaitAny create a new `Listener` instance,
// register it and returns its channel to `range`
// through it.
func (e *Emitter) AwaitAny() chan *Event {
	if e.IsClosed() {
		return nil
	}

	lst := e.AddListener()
	return lst.c
}

// AddListener creates a new Listener
// instance, register it through `addListenerAction`
// and finally returns it.
func (e *Emitter) AddListener() *Listener {
	if e.IsClosed() {
		return nil
	}

	lst := Listener{
		c:          make(chan *Event),
		e:          e,
		closedLock: &sync.Mutex{},
	}
	actionsOnListeners <- addListenerAction(&lst, e)
	return &lst
}
