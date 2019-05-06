package logger

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDefaultManager(t *testing.T) {
	m := defaultManager.(*DefaultManager)
	l := NewSliceLogger().(*SliceLogger)
	require.NoError(t, Add(l))
	require.Len(t, m.loggers, 1)

	defer func() {
		require.NoError(t, Remove(l.ID()))
		require.Len(t, m.loggers, 0)
	}()

	t.Run("Global Data", func(t *testing.T) {
		AddGlobalData("key", "value")
		require.Len(t, m.globals, 1)

		RemoveGlobalData("key")
		require.Len(t, m.globals, 0)
	})

	t.Run("Tag", func(t *testing.T) {
		SetTag("tag")
		defer SetTag("")
		require.Equal(t, "tag", m.tag)
		assert.Equal(t, "tag", Tag())
		assert.Equal(t, "tag", FullTag())
	})

	t.Run("ID", func(t *testing.T) {
		assert.NotEmpty(t, ID())
	})

	t.Run("Error", func(t *testing.T) {
		defer l.clear()
		Error("a", "b")

		require.Len(t, l.data, 1, "no logs added")
		require.Equal(t, "[ERROR]a b\n", l.data[0])
	})

	t.Run("Errorf", func(t *testing.T) {
		defer l.clear()
		Errorf("%s %s", "a", "b")

		require.Len(t, l.data, 1, "no logs added")
		require.Equal(t, "[ERROR]a b\n", l.data[0])
	})

	t.Run("Info", func(t *testing.T) {
		defer l.clear()
		Info("a", "b")

		require.Len(t, l.data, 1, "no logs added")
		require.Equal(t, "[INFO]a b\n", l.data[0])
	})

	t.Run("Infof", func(t *testing.T) {
		defer l.clear()
		Infof("%s %s", "a", "b")

		require.Len(t, l.data, 1, "no logs added")
		require.Equal(t, "[INFO]a b\n", l.data[0])
	})

	t.Run("Debug", func(t *testing.T) {
		defer l.clear()
		Debug("a", "b")

		require.Len(t, l.data, 1, "no logs added")
		require.Equal(t, "[DEBUG]a b\n", l.data[0])
	})

	t.Run("Debugf", func(t *testing.T) {
		defer l.clear()
		Debugf("%s %s", "a", "b")

		require.Len(t, l.data, 1, "no logs added")
		require.Equal(t, "[DEBUG]a b\n", l.data[0])
	})

	t.Run("Log", func(t *testing.T) {
		defer l.clear()
		Log("a", "b")

		require.Len(t, l.data, 1, "no logs added")
		require.Equal(t, "a b\n", l.data[0])
	})

	t.Run("Logf", func(t *testing.T) {
		defer l.clear()
		Logf("%s %s", "a", "b")

		require.Len(t, l.data, 1, "no logs added")
		require.Equal(t, "a b\n", l.data[0])
	})
}

func TestDefaultManagerClose(t *testing.T) {
	m := defaultManager.(*DefaultManager)
	l := NewSliceLogger().(*SliceLogger)
	require.NoError(t, Add(l))
	require.Len(t, m.loggers, 1)

	require.Len(t, Close(), 0)
	require.Len(t, m.loggers, 0)
}

func TestDefaultManagerSub(t *testing.T) {
	m := defaultManager.(*DefaultManager)
	SetTag("[parent]")
	l := NewSliceLogger().(*SliceLogger)
	require.NoError(t, Add(l))
	require.Len(t, m.loggers, 1)

	defer func() {
		require.Len(t, Close(), 0)
		assert.Len(t, m.loggers, 0)
		assert.Len(t, m.children, 0)
	}()

	sm := NewSubManager("[child]")
	assert.Equal(t, "[parent][child]", sm.FullTag())
	assert.Equal(t, "[parent]", FullTag())

	sm.Log("foo")

	require.Len(t, l.data, 1)
	assert.Equal(t, "[parent][child] foo\n", l.data[0])
}
