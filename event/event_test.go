package event

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type Field struct {
	Changed *Emitter
	Ready   *Emitter
}

func newTestSource() *Field {
	source := Field{}
	InitSource(&source, &source.Changed, &source.Ready)
	go func() {
		for i := 0; i < 10; i++ {
			time.Sleep(20 * time.Millisecond)
			source.Ready.Invoke(i)
		}
		CloseEmitters(&source.Changed, &source.Ready)
	}()
	return &source
}

func TestEvent(t *testing.T) {
	t.Run("Invoke", func(t *testing.T) {
		t.Run("AwaitAny", func(t *testing.T) {
			source := newTestSource()
			counter := 0
			for e := range source.Ready.AwaitAny() {
				assert.NotNil(t, e)
				assert.Equal(t, e.Source, source)
				assert.Equal(t, e.Payload, counter)
				counter++
			}
			assert.Equal(t, 10, counter)
		})

		t.Run("AwaitOne", func(t *testing.T) {
			source := newTestSource()

			for counter := 0; counter < 10; counter++ {
				e := source.Ready.AwaitOne()
				assert.NotNil(t, e)
				assert.Equal(t, e.Source, source)
				assert.Equal(t, e.Payload, counter)
			}
		})

		t.Run("Listen", func(t *testing.T) {
			source := newTestSource()
			counter := 0

			source.Ready.Listen(func(e *Event) {
				assert.NotNil(t, e)
				assert.Equal(t, e.Source, source)
				assert.Equal(t, e.Payload, counter)
				counter++
			})

			time.Sleep(500 * time.Millisecond)
			assert.Equal(t, 10, counter)
		})
	})

	t.Run("InitSource", func(t *testing.T) {
		source := Field{}
		assert.Nil(t, source.Changed)
		assert.Nil(t, source.Ready)
		InitSource(source, &source.Changed, &source.Ready)
		assert.NotNil(t, source.Changed)
		assert.NotNil(t, source.Ready)
	})

	t.Run("Emitter", func(t *testing.T) {
		t.Run("NewEmitter", func(t *testing.T) {
			source := Field{}
			source.Changed = NewEmitter(&source)

			assert.NotNil(t, source.Changed)
			assert.Equal(t, &source, source.Changed.source)
			assert.NotNil(t, source.Changed.listeners)
			assert.Equal(t, 0, len(source.Changed.listeners))
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
			source.Changed.Close()
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
