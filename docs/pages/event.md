{{ useLayout(".layout.njk") }}
{{ title("CIMA virtual-server") }}
{{ subtitle("event package") }}

# [virtual-server](./index) ⟶ {{ meta.subtitle }}



Package event implements an event framework that allows events to be listened
and emitted from multiple goroutins.

## Table of contents

[[TOC]]

## Example

```go

package main import "event"

type ExampleSource struct {

    SomeThingChanged event.Emitter

}

func main() {

    s := ExampleSource {}

    // create an emitter instance
    s.AnEvent = event.NewEmitter(s)

    // range over all emitted events
    go func() {
      for ev := range s.AnEvent.AwaitAny() {
         fmt.Println(ev.Source, ev.Payload)
      }
    }()

    // await for a single event emission
    go func() {
      ev, valid := s.AnEvent.AwaitOne();
      if valid {
         fmt.Println(ev.Source, ev.Payload)
      }
    }()

    hndl := func(ev *event.Event) {
      fmt.Println(ev.Source, ev.Payload)
    }

    // register a function that will
    // executed on each event emission
    s.AnEvent.Listen(hndl)

}

```

## Usage

#### type Emitter

```go
type Emitter struct {
}
```

Emitter is an object which can emits multiple events of the same kind,

Event emission can be listened by multiple listeners, and can happens from
multiple goroutines. Event invocation also is completely synchronized and can
occur in multiple goroutines.

#### func  NewEmitter

```go
func NewEmitter(source Source) *Emitter
```
NewEmitter return a new instance of Emitter, linked `source` argument.

#### func (*Emitter) AddListener

```go
func (e *Emitter) AddListener() *Listener
```
AddListener creates a new Listener instance, register it through
`addListenerAction` and finally returns it.

#### func (*Emitter) AwaitAny

```go
func (e *Emitter) AwaitAny() chan *Event
```
AwaitAny create a new `Listener` instance, register it and returns its channel
to `range` through it.

#### func (*Emitter) AwaitOne

```go
func (e *Emitter) AwaitOne() *Event
```
AwaitOne registers a new `Listener` on an Emitter, then waits for an event to
occurs on it, and finally unregisters the listener instance.

#### func (*Emitter) Clear

```go
func (e *Emitter) Clear()
```
Clear removes all listener of the `Emitter`

#### func (*Emitter) Count

```go
func (e *Emitter) Count() int
```
Count returns total number of listeners currently listening on this emitter.

#### func (*Emitter) Invoke

```go
func (e *Emitter) Invoke(payload interface{})
```
Invoke causes all listeners of this `Emitter` to receive an `Event` instance
with `Emitter` source as source, and `payload` argument as `payload`

#### func (*Emitter) Listen

```go
func (e *Emitter) Listen(fn Handler) *Listener
```
Listen add a listener to the emitter that executes a function for each event
emitted.

#### func (*Emitter) Stop

```go
func (e *Emitter) Stop()
```
Stop all listeners and closes the emitter.

#### type Event

```go
type Event struct {
	Source  Source
	Payload interface{}
}
```

Event is a struct that contains infgormations about a single event emitted. It
contains two fields:

`Source`, which contains the Object which calls `Invoke` method, and `Payload`
whic contains informations specifics to the event type.

#### type Handler

```go
type Handler func(ev *Event)
```

Handler ...

#### type Listener

```go
type Listener struct {
}
```

Listener is a single listener which is listening on the event.

#### func (*Listener) Stop

```go
func (l *Listener) Stop()
```
Stop listening new events invoked on a single listener.

#### type Source

```go
type Source interface{}
```

Source is an interface that represents any object which emits one or more
events.
