package token_test

import (
	"testing"
	"time"

	"github.com/rainb0w-clwn/go_auth_limiter/internal/bucket/refillrate"
	"github.com/rainb0w-clwn/go_auth_limiter/internal/bucket/token"
	"github.com/stretchr/testify/require"
)

func TestNewBucket(t *testing.T) {
	refill := refillrate.New(1, time.Second)
	b := token.New(10, refill)

	require.Equal(t, 10, b.GetSize())
	require.Equal(t, 10, b.GetTokenCount())
	require.True(t, b.Full())
}

func TestGetToken(t *testing.T) {
	refill := refillrate.New(1, time.Second)
	b := token.New(10, refill)

	b.GetToken(3)

	require.Equal(t, 7, b.GetTokenCount())
	require.False(t, b.Full())
}

func TestReset(t *testing.T) {
	refill := refillrate.New(1, time.Second)
	b := token.New(10, refill)

	b.GetToken(5)
	require.Equal(t, 5, b.GetTokenCount())

	b.Reset()

	require.Equal(t, 10, b.GetTokenCount())
	require.True(t, b.Full())
}

func TestRefillAddsTokens(t *testing.T) {
	refill := refillrate.New(2, time.Second) // 2 токена в секунду
	b := token.New(10, refill)

	b.GetToken(6)
	require.Equal(t, 4, b.GetTokenCount())

	time.Sleep(1100 * time.Millisecond)
	b.Refill()

	// за ~1 секунду должно добавиться минимум 2 токена
	require.GreaterOrEqual(t, b.GetTokenCount(), 6)
	require.LessOrEqual(t, b.GetTokenCount(), 10)
}

func TestRefillDoesNotOverflow(t *testing.T) {
	refill := refillrate.New(100, time.Second)
	b := token.New(10, refill)

	b.GetToken(1)
	require.Equal(t, 9, b.GetTokenCount())

	time.Sleep(1100 * time.Millisecond)
	b.Refill()

	require.Equal(t, 10, b.GetTokenCount())
	require.True(t, b.Full())
}

func TestLastRefillUpdatedOnlyWhenTokensAdded(t *testing.T) {
	refill := refillrate.New(1, time.Second)
	b := token.New(10, refill)

	last := b.GetLastRefill()

	// слишком мало времени — refill не должен сработать
	time.Sleep(100 * time.Millisecond)
	b.Refill()

	require.Equal(t, last, b.GetLastRefill())

	// теперь refill должен добавить токены и обновить время
	b.GetToken(1)
	time.Sleep(1100 * time.Millisecond)
	b.Refill()

	require.True(t, b.GetLastRefill().After(last))
}

func TestFull(t *testing.T) {
	refill := refillrate.New(1, time.Second)
	b := token.New(5, refill)

	require.True(t, b.Full())

	b.GetToken(1)
	require.False(t, b.Full())
}
