package refillrate_test

import (
	"testing"
	"time"

	"github.com/rainb0w-clwn/go_auth_limiter/internal/bucket/refillrate"
	"github.com/stretchr/testify/require"
)

func TestNewRefillRate(t *testing.T) {
	count := 10
	duration := 5 * time.Second

	rate := refillrate.New(count, duration)

	require.Equal(t, count, rate.GetCount())
	require.Equal(t, duration, rate.GetTime())
}

func TestGetCount(t *testing.T) {
	rate := refillrate.New(3, time.Minute)

	result := rate.GetCount()

	require.Equal(t, 3, result)
}

func TestGetTime(t *testing.T) {
	duration := 2 * time.Hour
	rate := refillrate.New(1, duration)

	result := rate.GetTime()

	require.Equal(t, duration, result)
}

func TestGetTTL(t *testing.T) {
	duration := 2 * time.Hour
	rate := refillrate.New(1, duration)

	result := rate.GetTime()

	require.Equal(t, duration, result)
}

func TestZeroValues(t *testing.T) {
	rate := refillrate.New(0, 0)

	require.Equal(t, 0, rate.GetCount())
	require.Equal(t, time.Duration(0), rate.GetTime())
}
