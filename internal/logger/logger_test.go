package logger

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLogger(t *testing.T) {
	const msg = "log msg"
	t.Run("log with exact level", func(t *testing.T) {
		out := &bytes.Buffer{}
		logg := New(Warn, out)
		logg.Warn(msg)
		require.Contains(t, out.String(), msg)
	})
	t.Run("log with higher level", func(t *testing.T) {
		out := &bytes.Buffer{}
		log := New(Info, out)
		log.Warn(msg)
		require.Contains(t, out.String(), msg)
	})
	t.Run("log with lower level", func(t *testing.T) {
		out := &bytes.Buffer{}
		log := New(Error, out)
		log.Warn(msg)
		require.Empty(t, out.String())
	})
	t.Run("get level from string", func(t *testing.T) {
		for level := range LevelMap {
			require.NotPanics(t, func() { GetLevelOrPanic(level) })
			require.NotPanics(t, func() { GetLevelOrPanic(strings.ToUpper(level)) })
		}
		require.PanicsWithError(t, ErrUnknownLevel.Error(), func() { GetLevelOrPanic("unknown") })
	})
}
