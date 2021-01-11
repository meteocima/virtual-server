package event

// internal struct used to
// pass actions info, through an `Emitter.actionsOnListeners`
// channel, between requesting goroutines
// and the `Emitter` internal one.
type listenerAction struct {
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
	listenerActionClear
	// action used to emit a new event
	listenerActionEmit
	// action used to query for listeners count
	listenerActionCount
)

// functions below create listenerAction
// for of any possible kind of action.

func countListenersAction() listenerAction {
	countResp := make(chan int)
	return listenerAction{
		kind:      listenerActionCount,
		countResp: countResp,
	}
}

func addListenerAction(listener *Listener) listenerAction {
	return listenerAction{
		kind:     listenerActionAdd,
		listener: listener,
	}
}

func removeListenerAction(listener *Listener) listenerAction {
	return listenerAction{
		kind:     listenerActionRemove,
		listener: listener,
	}
}

func clearListenersAction() listenerAction {
	return listenerAction{
		kind: listenerActionClear,
	}
}

func emitListenerAction(event *Event) listenerAction {
	return listenerAction{
		kind:  listenerActionEmit,
		event: event,
	}
}

// loop though all listener actions emitted
// from `actionsOnListeners` channel, and
// execute the action specific for any of them.
func (e *Emitter) executeListenerActions() {
	for action := range e.actionsOnListeners {
		switch action.kind {
		case listenerActionAdd:
			e.listeners[action.listener] = struct{}{}
		case listenerActionClear:
			e.listeners = map[*Listener]struct{}{}
		case listenerActionEmit:
			/* YODO */
		case listenerActionRemove:
			delete(e.listeners, action.listener)
		case listenerActionCount:
			action.countResp <- len(e.listeners)
			close(action.countResp)
		default:
			panic("Unknown action kind")
		}
	}
}
