package gb_test

import (
	"testing"
	"time"

	"github.com/rainb0w-clwn/go_auth_limiter/internal/bucket/refillrate"
	"github.com/rainb0w-clwn/go_auth_limiter/internal/limiter"
	"github.com/rainb0w-clwn/go_auth_limiter/internal/limiter/composite"
	limitermocks "github.com/rainb0w-clwn/go_auth_limiter/internal/limiter/mocks"
	"github.com/rainb0w-clwn/go_auth_limiter/internal/limiter/tokenbucket/gb"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestTokenBucketGB_Sweep(t *testing.T) {
	size := 3
	ttl := time.Millisecond * 100
	refillRate := refillrate.New(1, time.Hour*1) // "disable" auto refill with long rate
	tokenBucketLimiter := getTokenBucketLimiter(t, size, refillRate)
	gb := gb.New(tokenBucketLimiter, ttl)

	// init token buckets
	tokenBucketLimiter.SetRequestCost(size)
	identities := []limiter.UserIdentityDto{
		{limiter.IPLimit.String(): "192.168.1.1"},
		{limiter.LoginLimit.String(): "root"},
		{limiter.LoginLimit.String(): "abc"},
		{limiter.PasswordLimit.String(): "123456"},
	}
	for _, identity := range identities {
		tokenBucketLimiter.SatisfyLimit(identity)
		tokenBucketLimiter.ResetLimit(identity)
	}
	require.Len(t, tokenBucketLimiter.GetBuckets(), len(identities))

	// drain one identity to check GB sweeping only full buckets (despite ttl expired)
	tokenBucketLimiter.SatisfyLimit(identities[2])

	// sleep to make buckets refill date outdated
	time.Sleep(ttl)

	// run db
	gb.Sweep()

	require.Len(t, tokenBucketLimiter.GetBuckets(), 1)
}

func getMockLimitStorage(t *testing.T, types []limiter.Type, values []int) *limitermocks.MockIStorage {
	t.Helper()

	mockLimits := make(limiter.Limits, 0, len(types))
	for i, limitType := range types {
		mockLimits = append(mockLimits, limiter.Limit{
			LimitType: limitType,
			Value:     values[i],
		})
	}

	limitStorage := limitermocks.NewMockIStorage(t)
	limitStorage.EXPECT().GetLimitsByTypes(mock.AnythingOfType("[]string")).Return(&mockLimits, nil)

	return limitStorage
}

func getTokenBucketLimiter(t *testing.T, size int, refillRate refillrate.RefillRate) *composite.Limiter {
	t.Helper()

	return composite.New(
		getMockLimitStorage(
			t,
			[]limiter.Type{limiter.LoginLimit, limiter.IPLimit, limiter.PasswordLimit},
			[]int{size, size, size},
		),
		refillRate,
	)
}
