package logger

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestManagerID(t *testing.T) {
	t.Parallel()

	t.Run("id gets returned by ID()", func(t *testing.T) {
		t.Parallel()
		m := &DefaultManager{id: "id"}
		assert.Equal(t, "id", m.ID())
	})

	t.Run("NewManager() sets an ID", func(t *testing.T) {
		t.Parallel()
		assert.NotEmpty(t, NewManager().ID())
	})
}

func TestManagerTag(t *testing.T) {
	t.Parallel()

	t.Run("tag gets returned by Tag()", func(t *testing.T) {
		t.Parallel()
		m := &DefaultManager{tag: "tag"}
		assert.Equal(t, "tag", m.Tag())
	})

	t.Run("tag gets returned by FullTag()", func(t *testing.T) {
		t.Parallel()
		m := &DefaultManager{tag: "tag"}
		assert.Equal(t, "tag", m.FullTag())
	})

	t.Run("parents tag gets returned by FullTag()", func(t *testing.T) {
		t.Parallel()
		m := NewManagerWithTag("grand-parent")
		sm := m.NewSubManager("parent")
		ssm := sm.NewSubManager("child")

		expected := fmt.Sprintf("%s%s%s", m.Tag(), sm.Tag(), ssm.Tag())
		assert.Equal(t, expected, ssm.FullTag())

		expected = fmt.Sprintf("%s%s", m.Tag(), sm.Tag())
		assert.Equal(t, expected, sm.FullTag())

		expected = m.Tag()
		assert.Equal(t, expected, m.FullTag())
	})

	t.Run("tag gets set by SetTag()", func(t *testing.T) {
		t.Parallel()
		m := NewManager()
		m.SetTag("new-tag")
		assert.NotEmpty(t, "new-tag", m.Tag())
	})

	t.Run("NewManagerWithTag() sets a tag", func(t *testing.T) {
		t.Parallel()
		assert.Equal(t, "tag", NewManagerWithTag("tag").Tag())
	})
}

func TestManagerGlobalData(t *testing.T) {
	t.Parallel()

	nm := NewManager()
	m := nm.(*DefaultManager)
	m.AddGlobalData("key", "value")
	require.Len(t, m.globals, 1, "no data added")

	m.RemoveGlobalData("key")
	require.Empty(t, m.globals, "no data removed")
}

func TestManagerSubManager(t *testing.T) {
	t.Parallel()

	nm := NewManager()
	m := nm.(*DefaultManager)
	sb := m.NewSubManager("child")
	require.Len(t, m.children, 1, "no manager added")
	require.Equal(t, m, sb.(*DefaultManager).parent, "wrong parent")

	sb.Close()
	require.Empty(t, m.children, "no manager removed")
}

func TestManagerLogger(t *testing.T) {
	t.Parallel()

	t.Run("Add loggers", func(t *testing.T) {
		t.Parallel()
		nm := NewManager()
		m := nm.(*DefaultManager)

		l := NewSliceLogger()
		require.NoError(t, m.Add(l))

		l2 := NewSliceLogger()
		l2.(*SliceLogger).id = "fake-id"
		require.NoError(t, m.Add(l2))

		require.Len(t, m.loggers, 2, "no loggers added")
	})

	t.Run("Add duplicates", func(t *testing.T) {
		t.Parallel()
		nm := NewManager()
		m := nm.(*DefaultManager)

		l := NewSliceLogger()
		require.NoError(t, m.Add(l))

		err := m.Add(l)
		require.Error(t, err)
		assert.Equal(t, ErrAlreadyExist, err)
	})

	t.Run("Remove a logger", func(t *testing.T) {
		t.Parallel()
		nm := NewManager()
		m := nm.(*DefaultManager)

		l := NewSliceLogger()
		require.NoError(t, m.Add(l))

		l2 := NewSliceLogger()
		l2.(*SliceLogger).id = "fake-id"
		require.NoError(t, m.Add(l2))

		require.Len(t, m.loggers, 2, "no loggers added")

		err := m.Remove(l.ID())
		require.NoError(t, err)
		require.Len(t, m.loggers, 1, "wrong logger removed")
	})

	t.Run("Close loggers", func(t *testing.T) {
		t.Parallel()
		nm := NewManager()
		m := nm.(*DefaultManager)

		l := NewSliceLogger()
		require.NoError(t, m.Add(l))

		l2 := NewSliceLogger()
		l2.(*SliceLogger).id = "fake-id"
		require.NoError(t, m.Add(l2))

		require.Len(t, m.loggers, 2, "no loggers added")

		err := m.Close()
		require.Empty(t, err)
		require.Empty(t, m.loggers, "no loggers removed")
	})
}

func TestManagerLog(t *testing.T) {
	t.Parallel()

	t.Run("Error", func(t *testing.T) {
		t.Parallel()
		m := NewManager()

		lo := NewSliceLogger()
		l := lo.(*SliceLogger)
		require.NoError(t, m.Add(l))

		m.Error("a", "b")

		require.Len(t, l.data, 1, "no logs added")
		require.Equal(t, "[ERROR]a b\n", l.data[0])
	})

	t.Run("Errorf", func(t *testing.T) {
		t.Parallel()
		m := NewManager()

		lo := NewSliceLogger()
		l := lo.(*SliceLogger)
		require.NoError(t, m.Add(l))

		m.Errorf("%s %s", "a", "b")

		require.Len(t, l.data, 1, "no logs added")
		require.Equal(t, "[ERROR]a b\n", l.data[0])
	})

	t.Run("Info", func(t *testing.T) {
		t.Parallel()
		m := NewManager()

		lo := NewSliceLogger()
		l := lo.(*SliceLogger)
		require.NoError(t, m.Add(l))

		m.Info("a", "b")

		require.Len(t, l.data, 1, "no logs added")
		require.Equal(t, "[INFO]a b\n", l.data[0])
	})

	t.Run("Infof", func(t *testing.T) {
		t.Parallel()
		m := NewManager()

		lo := NewSliceLogger()
		l := lo.(*SliceLogger)
		require.NoError(t, m.Add(l))

		m.Infof("%s %s", "a", "b")

		require.Len(t, l.data, 1, "no logs added")
		require.Equal(t, "[INFO]a b\n", l.data[0])
	})

	t.Run("Debug", func(t *testing.T) {
		t.Parallel()
		m := NewManager()

		lo := NewSliceLogger()
		l := lo.(*SliceLogger)
		require.NoError(t, m.Add(l))

		m.Debug("a", "b")

		require.Len(t, l.data, 1, "no logs added")
		require.Equal(t, "[DEBUG]a b\n", l.data[0])
	})

	t.Run("Debugf", func(t *testing.T) {
		t.Parallel()
		m := NewManager()

		lo := NewSliceLogger()
		l := lo.(*SliceLogger)
		require.NoError(t, m.Add(l))

		m.Debugf("%s %s", "a", "b")

		require.Len(t, l.data, 1, "no logs added")
		require.Equal(t, "[DEBUG]a b\n", l.data[0])
	})

	t.Run("Log", func(t *testing.T) {
		t.Parallel()
		m := NewManager()

		lo := NewSliceLogger()
		l := lo.(*SliceLogger)
		require.NoError(t, m.Add(l))

		m.Log("a", "b")

		require.Len(t, l.data, 1, "no logs added")
		require.Equal(t, "a b\n", l.data[0])
	})

	t.Run("Logf", func(t *testing.T) {
		t.Parallel()
		m := NewManager()

		lo := NewSliceLogger()
		l := lo.(*SliceLogger)
		require.NoError(t, m.Add(l))

		m.Logf("%s %s", "a", "b")

		require.Len(t, l.data, 1, "no logs added")
		require.Equal(t, "a b\n", l.data[0])
	})

	t.Run("Info with parents", func(t *testing.T) {
		t.Parallel()

		m := NewManager()
		lo1 := NewSliceLogger()
		l1 := lo1.(*SliceLogger)
		require.NoError(t, m.Add(l1))

		sm := m.NewSubManager("[child]")
		lo2 := NewSliceLogger()
		l2 := lo2.(*SliceLogger)
		l2.id = "fake-id"
		require.NoError(t, sm.Add(l2))

		sm.Info("a", "b")

		require.Len(t, l1.data, 1, "no logs added")
		assert.Equal(t, "[INFO][child] a b\n", l1.data[0])

		require.Len(t, l2.data, 1, "no logs added")
		assert.Equal(t, "[INFO][child] a b\n", l2.data[0])
	})

	t.Run("Error with parents", func(t *testing.T) {
		t.Parallel()

		m := NewManager()
		lo1 := NewSliceLogger()
		l1 := lo1.(*SliceLogger)
		require.NoError(t, m.Add(l1))

		sm := m.NewSubManager("[child]")
		lo2 := NewSliceLogger()
		l2 := lo2.(*SliceLogger)
		l2.id = "fake-id"
		require.NoError(t, sm.Add(l2))

		sm.Error("a", "b")

		require.Len(t, l1.data, 1, "no logs added")
		assert.Equal(t, "[ERROR][child] a b\n", l1.data[0])

		require.Len(t, l2.data, 1, "no logs added")
		assert.Equal(t, "[ERROR][child] a b\n", l2.data[0])
	})

	t.Run("Debug with parents", func(t *testing.T) {
		t.Parallel()

		m := NewManager()
		lo1 := NewSliceLogger()
		l1 := lo1.(*SliceLogger)
		require.NoError(t, m.Add(l1))

		sm := m.NewSubManager("[child]")
		lo2 := NewSliceLogger()
		l2 := lo2.(*SliceLogger)
		l2.id = "fake-id"
		require.NoError(t, sm.Add(l2))

		sm.Debug("a", "b")

		require.Len(t, l1.data, 1, "no logs added")
		assert.Equal(t, "[DEBUG][child] a b\n", l1.data[0])

		require.Len(t, l2.data, 1, "no logs added")
		assert.Equal(t, "[DEBUG][child] a b\n", l2.data[0])
	})

	t.Run("Log with tagged parents", func(t *testing.T) {
		t.Parallel()

		m := NewManagerWithTag("[parent]")
		lo1 := NewSliceLogger()
		l1 := lo1.(*SliceLogger)
		require.NoError(t, m.Add(l1))

		sm := m.NewSubManager("[child]")
		lo2 := NewSliceLogger()
		l2 := lo2.(*SliceLogger)
		l2.id = "fake-id"
		require.NoError(t, sm.Add(l2))

		sm.Log("a", "b")

		require.Len(t, l1.data, 1, "no logs added")
		assert.Equal(t, "[parent][child] a b\n", l1.data[0])

		require.Len(t, l2.data, 1, "no logs added")
		assert.Equal(t, "[parent][child] a b\n", l2.data[0])

		m.Log("c", "d")

		require.Len(t, l1.data, 2, "no logs added")
		assert.Equal(t, "[parent] c d\n", l1.data[1])

		require.Len(t, l2.data, 1, "no logs should have been added")
	})
}
