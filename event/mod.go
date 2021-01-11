// Package event implements an event framework
// that allows events to be listened and emitted from
// multiple goroutins.
//
// ## Table of contents
//
// [[TOC]]
//
// ## Example
//
// ```go
//
// package main
// import "event"
//
// type ExampleSource struct {
//    SomeThingChanged event.Emitter
// }
//
// func main() {
//  s := ExampleSource {}
//
//  // create an emitter instance
//  s.AnEvent = event.NewEmitter(s)
//
//  // range over all emitted events
//  go func() {
//    for ev := range s.AnEvent.AwaitAny() {
//       fmt.Println(ev.Source, ev.Payload)
//    }
//  }()
//
//  // await for a single event emission
//  go func() {
//    ev, valid := s.AnEvent.AwaitOne();
//    if valid {
//       fmt.Println(ev.Source, ev.Payload)
//    }
//  }()
//
//  hndl := func(ev *event.Event) {
//    fmt.Println(ev.Source, ev.Payload)
//  }
//
//  // register a function that will
//  // executed on each event emission
//  s.AnEvent.Listen(hndl)
//
// }
//
// ```
//
package event
