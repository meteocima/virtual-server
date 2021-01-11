package event

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type Field struct {
	Changed *Emitter
}

func TestEvent(t *testing.T) {
	t.Run("Emitter", func(t *testing.T) {
		t.Run("NewEmitter", func(t *testing.T) {
			source := Field{}
			source.Changed = NewEmitter(&source)

			assert.NotNil(t, source.Changed)
			assert.Equal(t, &source, source.Changed.source)
			assert.NotNil(t, source.Changed.listeners)
			assert.NotNil(t, source.Changed.actionsOnListeners)
			assert.Equal(t, 0, len(source.Changed.listeners))
			assert.Equal(t, 0, len(source.Changed.actionsOnListeners))
		})

		t.Run("Listen", func(t *testing.T) {
			source := Field{}
			source.Changed = NewEmitter(&source)
			assert.Equal(t, 0, source.Changed.Count())
			source.Changed.Listen(func(e *Event) {})
			assert.Equal(t, 1, source.Changed.Count())
		})

		t.Run("Clear", func(t *testing.T) {
			source := Field{}
			source.Changed = NewEmitter(&source)
			source.Changed.Listen(func(e *Event) {})
			source.Changed.Listen(func(e *Event) {})
			assert.Equal(t, 2, source.Changed.Count())
			source.Changed.Clear()
			assert.Equal(t, 0, source.Changed.Count())
		})

		t.Run("Stop", func(t *testing.T) {
			source := Field{}
			source.Changed = NewEmitter(&source)

			listener := source.Changed.Listen(func(e *Event) {})

			assert.Equal(t, 1, source.Changed.Count())
			listener.Stop()
			assert.Equal(t, 0, source.Changed.Count())
		})

	})
}
