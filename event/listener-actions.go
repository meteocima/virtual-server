package event

// internal struct used to
// pass actions info, through an `Emitter.actionsOnListeners`
// channel, between requesting goroutines
// and the `Emitter` internal one.
type listenerAction struct {
	emitter   *Emitter
	kind      actionKind
	listener  *Listener
	event     *Event
	countResp chan int
}

type actionKind int

const (
	// action used to add a new listener
	listenerActionAdd actionKind = iota
	// action used to remove an existing listener
	listenerActionRemove
	// action used to remove all listeners
	listenerActionCloseEmitter
	// action used to emit a new event
	listenerActionEmit
	// action used to query for listeners count
	listenerActionCount
)

// functions below create listenerAction
// for of any possible kind of action.

func countListenersAction(emitter *Emitter) listenerAction {
	countResp := make(chan int)
	return listenerAction{
		emitter:   emitter,
		kind:      listenerActionCount,
		countResp: countResp,
	}
}

func addListenerAction(listener *Listener, emitter *Emitter) listenerAction {
	return listenerAction{
		emitter:  emitter,
		kind:     listenerActionAdd,
		listener: listener,
	}
}

func removeListenerAction(listener *Listener, emitter *Emitter) listenerAction {
	return listenerAction{
		emitter:  emitter,
		kind:     listenerActionRemove,
		listener: listener,
	}
}

func closeEmitterListenersAction(emitter *Emitter) listenerAction {
	return listenerAction{
		emitter: emitter,
		kind:    listenerActionCloseEmitter,
	}
}

func emitListenerAction(event *Event, emitter *Emitter) listenerAction {
	countResp := make(chan int)

	return listenerAction{
		emitter:   emitter,
		kind:      listenerActionEmit,
		event:     event,
		countResp: countResp,
	}
}

var actionsOnListeners = make(chan listenerAction)

// loop though all listener actions emitted
// from `actionsOnListeners` channel, and
// execute the action specific for any of them.
func init() {
	go func() {
		for action := range actionsOnListeners {
			switch action.kind {
			case listenerActionAdd:
				action.emitter.listeners[action.listener] = struct{}{}
			case listenerActionCloseEmitter:
				for l := range action.emitter.listeners {
					l.killChannel()
				}
				action.emitter.listeners = map[*Listener]struct{}{}
			case listenerActionEmit:
				for l := range action.emitter.listeners {
					l.c <- action.event
				}
				action.countResp <- 0
			case listenerActionRemove:
				action.listener.killChannel()
				delete(action.emitter.listeners, action.listener)
			case listenerActionCount:
				action.countResp <- len(action.emitter.listeners)
				close(action.countResp)
			default:
				panic("Unknown action kind")
			}
		}
	}()
}
