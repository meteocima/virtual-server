package event

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
// Event emission can be listened by multiple listeners,
// and can happens from multiple goroutines. Event invocation also
// is completely synchronized and can occur in multiple goroutines.
type Emitter struct {
	source             Source
	listeners          map[*Listener]struct{}
	actionsOnListeners chan listenerAction
}

// NewEmitter return a new instance
// of Emitter, linked `source` argument.
func NewEmitter(source Source) *Emitter {
	emitter := &Emitter{
		source:             source,
		listeners:          map[*Listener]struct{}{},
		actionsOnListeners: make(chan listenerAction),
	}
	go emitter.executeListenerActions()
	return emitter
}

// Listener is a single listener
// which is listening on the event.
type Listener struct {
	c chan *Event
	e *Emitter
}

// Stop listening new events invoked
// on a single listener.
func (l *Listener) Stop() {
	close(l.c)
	l.e.actionsOnListeners <- removeListenerAction(l)
}

// Stop all listeners and closes
// the emitter.
func (e *Emitter) Stop() {
	//e.actionsOnListeners <- removeListenerAction(&l)
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
	e.actionsOnListeners <- emitListenerAction(&event)
}

// Handler ...
type Handler func(ev *Event)

// Count returns total number of
// listeners currently listening
// on this emitter.
func (e *Emitter) Count() int {
	action := countListenersAction()
	e.actionsOnListeners <- action
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

// Clear removes all listener of
// the `Emitter`
func (e *Emitter) Clear() {
	e.actionsOnListeners <- clearListenersAction()
}

// AwaitOne registers a new `Listener` on an
// Emitter, then waits for an event to occurs
// on it, and finally unregisters the listener
// instance.
func (e *Emitter) AwaitOne() *Event {
	lst := e.AddListener()
	event := <-lst.c
	e.actionsOnListeners <- removeListenerAction(lst)
	return event
}

// AwaitAny create a new `Listener` instance,
// register it and returns its channel to `range`
// through it.
func (e *Emitter) AwaitAny() chan *Event {
	lst := e.AddListener()
	return lst.c
}

// AddListener creates a new Listener
// instance, register it through `addListenerAction`
// and finally returns it.
func (e *Emitter) AddListener() *Listener {
	lst := Listener{make(chan *Event), e}
	e.actionsOnListeners <- addListenerAction(&lst)
	return &lst
}
