package tokenbucket_test

import (
	"testing"
	"time"

	"github.com/rainb0w-clwn/go_auth_limiter/internal/bucket/refillrate"
	"github.com/rainb0w-clwn/go_auth_limiter/internal/limiter"
	"github.com/rainb0w-clwn/go_auth_limiter/internal/limiter/tokenbucket"
	"github.com/stretchr/testify/require"
)

func TestTokenBucketLimiter_SatisfyLimit(t *testing.T) { //nolint:funlen
	bucketKey := "ip"
	identity := limiter.UserIdentityDto{bucketKey: "192.168.1.1"}

	t.Run("limit reached", func(t *testing.T) {
		bucketSize := 3
		refillRate := refillrate.New(3, time.Second*1)
		tokenBucketLimiter := tokenbucket.New(bucketKey, bucketSize, refillRate)

		// 3 requests allowed
		for i := 0; i < bucketSize; i++ {
			satisfies, err := tokenBucketLimiter.SatisfyLimit(identity)
			require.True(t, satisfies)
			require.NoError(t, err)
		}

		// 4th and following requests denied
		for i := 0; i < bucketSize; i++ {
			satisfies, err := tokenBucketLimiter.SatisfyLimit(identity)
			require.False(t, satisfies)
			require.NoError(t, err)
		}
	})

	t.Run("too expensive request", func(t *testing.T) {
		bucketSize := 3
		refillRate := refillrate.New(3, time.Second*1)
		tokenBucketLimiter := tokenbucket.New(bucketKey, bucketSize, refillRate)

		tokenBucketLimiter.SetRequestCost(bucketSize + 1)
		satisfies, err := tokenBucketLimiter.SatisfyLimit(identity)

		require.False(t, satisfies)
		require.NoError(t, err)
	})

	t.Run("simple refill", func(t *testing.T) {
		bucketSize := 3
		refillRate := refillrate.New(3, time.Second*1)
		tokenBucketLimiter := tokenbucket.New(bucketKey, bucketSize, refillRate)

		// 3 requests allowed
		for i := 0; i < bucketSize; i++ {
			satisfies, err := tokenBucketLimiter.SatisfyLimit(identity)
			require.True(t, satisfies)
			require.NoError(t, err)
		}

		// 4th denied
		satisfies, err := tokenBucketLimiter.SatisfyLimit(identity)
		require.False(t, satisfies)
		require.NoError(t, err)

		// after refill rate time we can make requests again
		time.Sleep(refillRate.GetTime())

		for i := 0; i < bucketSize; i++ {
			satisfies, err := tokenBucketLimiter.SatisfyLimit(identity)
			require.True(t, satisfies)
			require.NoError(t, err)
		}

		// 4th denied
		satisfies, err = tokenBucketLimiter.SatisfyLimit(identity)
		require.False(t, satisfies)
		require.NoError(t, err)
	})

	t.Run("partial refill", func(t *testing.T) {
		bucketSize := 3
		refillRate := refillrate.New(10, time.Second*1)
		tokenBucketLimiter := tokenbucket.New(bucketKey, bucketSize, refillRate)

		// 3 requests allowed
		for i := 0; i < bucketSize; i++ {
			satisfies, err := tokenBucketLimiter.SatisfyLimit(identity)
			require.True(t, satisfies)
			require.NoError(t, err)
		}

		// 4th denied
		satisfies, err := tokenBucketLimiter.SatisfyLimit(identity)
		require.False(t, satisfies)
		require.NoError(t, err)

		// refill rate = 10req/sec,
		// so we can make one more request after 0.1sec
		time.Sleep(time.Millisecond * 100)

		// one more request is allowed
		satisfies, err = tokenBucketLimiter.SatisfyLimit(identity)
		require.True(t, satisfies)
		require.NoError(t, err)

		// followings not
		satisfies, err = tokenBucketLimiter.SatisfyLimit(identity)
		require.False(t, satisfies)
		require.NoError(t, err)
	})

	t.Run("dynamic request cost", func(t *testing.T) {
		refillRate := refillrate.New(10, time.Second*1)
		tokenBucketLimiter := tokenbucket.New(bucketKey, 10, refillRate)

		// request with cost of 6 tokens after some time
		time.Sleep(time.Millisecond * 300)
		tokenBucketLimiter.SetRequestCost(6)
		satisfies, err := tokenBucketLimiter.SatisfyLimit(identity)
		require.True(t, satisfies)
		require.NoError(t, err)

		// check allowed requests remained (0 because only 4 tokens left)
		allowed, _ := tokenBucketLimiter.GetRequestsAllowed(identity)
		require.Zero(t, allowed)

		// wait 200 ms to refill and make request with cost of 5 tokens
		time.Sleep(time.Millisecond * 200)
		tokenBucketLimiter.SetRequestCost(5)
		satisfies, err = tokenBucketLimiter.SatisfyLimit(identity)
		require.True(t, satisfies)
		require.NoError(t, err)

		// check allowed requests remained (0 because only 1 token left)
		allowed, _ = tokenBucketLimiter.GetRequestsAllowed(identity)
		require.Zero(t, allowed)

		// wait for full refill and check requests allowed (2 request allowed with cost of 5 each)
		time.Sleep(time.Second * 1)
		allowed, _ = tokenBucketLimiter.GetRequestsAllowed(identity)
		require.Equal(t, 2, allowed)
	})

	t.Run("tricky refill rate #1", func(t *testing.T) {
		bucketSize := 3
		refillRate := refillrate.New(3, time.Second*3) // same as 1t/1sec
		tokenBucketLimiter := tokenbucket.New(bucketKey, bucketSize, refillRate)

		tokenBucketLimiter.SetRequestCost(3)
		_, _ = tokenBucketLimiter.SatisfyLimit(identity) // waste all tokens
		tokenBucketLimiter.SetRequestCost(1)

		// expect 1 token after 1 sec
		time.Sleep(time.Second * 1)
		allowed, _ := tokenBucketLimiter.GetRequestsAllowed(identity)
		require.Equal(t, 1, allowed)
	})

	t.Run("tricky refill rate #2", func(t *testing.T) {
		bucketSize := 3
		refillRate := refillrate.New(125, time.Second*150) // 125t/2.5min = same as 0.8(3)t/1sec
		tokenBucketLimiter := tokenbucket.New(bucketKey, bucketSize, refillRate)

		tokenBucketLimiter.SetRequestCost(3)
		_, _ = tokenBucketLimiter.SatisfyLimit(identity) // waste all tokens
		tokenBucketLimiter.SetRequestCost(1)

		// expect 1 full token after 2 sec
		time.Sleep(time.Second * 1)
		allowed, _ := tokenBucketLimiter.GetRequestsAllowed(identity)
		require.Equal(t, 0, allowed)

		time.Sleep(time.Second * 1)
		allowed, _ = tokenBucketLimiter.GetRequestsAllowed(identity)
		require.Equal(t, 1, allowed)
	})

	t.Run("multiple identity", func(t *testing.T) {
		bucketSize := 3
		refillRate := refillrate.New(3, time.Second*1)
		tokenBucketLimiter := tokenbucket.New(bucketKey, bucketSize, refillRate)

		// waste all tokens for first ip
		identity1 := limiter.UserIdentityDto{bucketKey: "192.168.1.1"}
		tokenBucketLimiter.SetRequestCost(3)
		satisfies, err := tokenBucketLimiter.SatisfyLimit(identity1)
		require.True(t, satisfies)
		require.NoError(t, err)

		// check no more allowed for first ip
		satisfies, err = tokenBucketLimiter.SatisfyLimit(identity1)
		require.False(t, satisfies)
		require.NoError(t, err)

		// check allowed for another ip
		identity2 := limiter.UserIdentityDto{bucketKey: "192.155.10.32"}
		satisfies, err = tokenBucketLimiter.SatisfyLimit(identity2)
		require.True(t, satisfies)
		require.NoError(t, err)

		// wait to refill, try one more
		time.Sleep(time.Second * 1)

		satisfies, err = tokenBucketLimiter.SatisfyLimit(identity1)
		require.True(t, satisfies)
		require.NoError(t, err)

		satisfies, err = tokenBucketLimiter.SatisfyLimit(identity2)
		require.True(t, satisfies)
		require.NoError(t, err)
	})
}

func TestTokenBucketLimiter_SatisfyLimit_Errors(t *testing.T) {
	t.Run("incorrect identity error", func(t *testing.T) {
		bucketSize := 3
		bucketKey := "ip"
		refillRate := refillrate.New(3, time.Second*1)
		tokenBucketLimiter := tokenbucket.New(bucketKey, bucketSize, refillRate)

		identity := limiter.UserIdentityDto{"login": "admin"}
		_, satisfyLimitErr := tokenBucketLimiter.SatisfyLimit(identity)
		require.ErrorIs(t, limiter.ErrIncorrectIdentity, satisfyLimitErr)

		_, getRequestsAllowedErr := tokenBucketLimiter.GetRequestsAllowed(identity)
		require.ErrorIs(t, limiter.ErrIncorrectIdentity, getRequestsAllowedErr)
	})
}

func TestTokenBucketLimiter_ResetLimit(t *testing.T) {
	bucketSize := 3
	bucketKey := "ip"
	refillRate := refillrate.New(1, time.Hour*1)
	identity := limiter.UserIdentityDto{bucketKey: "192.168.1.1"}

	t.Run("reset on empty bucket", func(t *testing.T) {
		tokenBucketLimiter := tokenbucket.New(bucketKey, bucketSize, refillRate)

		// drain bucket
		tokenBucketLimiter.SetRequestCost(bucketSize)
		tokenBucketLimiter.SatisfyLimit(identity)

		// check bucket "empty"
		satisfies, err := tokenBucketLimiter.SatisfyLimit(identity)
		require.NoError(t, err)
		require.False(t, satisfies)

		// reset & check again
		err = tokenBucketLimiter.ResetLimit(identity)
		require.NoError(t, err)

		satisfies, err = tokenBucketLimiter.SatisfyLimit(identity)
		require.NoError(t, err)
		require.True(t, satisfies)
	})

	t.Run("reset on full bucket", func(t *testing.T) {
		refillTime := time.Millisecond * 100
		tokenBucketLimiter := tokenbucket.New(
			bucketKey,
			bucketSize,
			refillrate.New(3, refillTime),
		)

		// drain bucket
		tokenBucketLimiter.SetRequestCost(bucketSize)
		tokenBucketLimiter.SatisfyLimit(identity)

		// wait for auto refill
		time.Sleep(refillTime)

		// check no error
		resetErr := tokenBucketLimiter.ResetLimit(identity)
		require.NoError(t, resetErr)

		// check reset don't increase size
		tokenBucketLimiter.SetRequestCost(bucketSize + 1)
		satisfies, err := tokenBucketLimiter.SatisfyLimit(identity)
		require.NoError(t, err)
		require.False(t, satisfies)

		// check normal
		tokenBucketLimiter.SetRequestCost(bucketSize)
		satisfies, err = tokenBucketLimiter.SatisfyLimit(identity)
		require.NoError(t, err)
		require.True(t, satisfies)
	})

	t.Run("reset on half full bucket", func(t *testing.T) {
		tokenBucketLimiter := tokenbucket.New(bucketKey, bucketSize, refillRate)

		// drain bucket
		tokenBucketLimiter.SatisfyLimit(identity)

		// try to make bucket size request on half full bucket
		tokenBucketLimiter.SetRequestCost(bucketSize)
		satisfies, err := tokenBucketLimiter.SatisfyLimit(identity)
		require.NoError(t, err)
		require.False(t, satisfies)

		// reset and check again
		tokenBucketLimiter.ResetLimit(identity)
		satisfies, err = tokenBucketLimiter.SatisfyLimit(identity)
		require.NoError(t, err)
		require.True(t, satisfies)
	})

	t.Run("cold reset", func(t *testing.T) {
		tokenBucketLimiter := tokenbucket.New(bucketKey, bucketSize, refillRate)

		// cold reset, on non-init bucket
		resetErr := tokenBucketLimiter.ResetLimit(identity)
		require.NoError(t, resetErr)

		tokenBucketLimiter.SetRequestCost(bucketSize)
		satisfies, err := tokenBucketLimiter.SatisfyLimit(identity)
		require.NoError(t, err)
		require.True(t, satisfies)
	})

	t.Run("multiple buckets", func(t *testing.T) {
		identity2 := limiter.UserIdentityDto{bucketKey: "10.25.13.3"}
		tokenBucketLimiter := tokenbucket.New(bucketKey, bucketSize, refillRate)
		tokenBucketLimiter.SetRequestCost(bucketSize)

		// drain buckets
		tokenBucketLimiter.SatisfyLimit(identity)
		tokenBucketLimiter.SatisfyLimit(identity2)

		// reset first bucket and check
		tokenBucketLimiter.ResetLimit(identity)

		satisfies, err := tokenBucketLimiter.SatisfyLimit(identity2)
		require.NoError(t, err)
		require.False(t, satisfies)

		satisfies, err = tokenBucketLimiter.SatisfyLimit(identity)
		require.NoError(t, err)
		require.True(t, satisfies)

		// reset second bucket and check
		tokenBucketLimiter.ResetLimit(identity2)

		satisfies, err = tokenBucketLimiter.SatisfyLimit(identity)
		require.NoError(t, err)
		require.False(t, satisfies)

		satisfies, err = tokenBucketLimiter.SatisfyLimit(identity2)
		require.NoError(t, err)
		require.True(t, satisfies)
	})
}

func TestTokenBucketLimiter_ResetLimit_Errors(t *testing.T) {
	t.Run("incorrect identity", func(t *testing.T) {
		bucketSize := 3
		bucketKey := "ip"
		refillRate := refillrate.New(1, time.Hour*1)
		tokenBucketLimiter := tokenbucket.New(bucketKey, bucketSize, refillRate)

		resetErr := tokenBucketLimiter.ResetLimit(limiter.UserIdentityDto{"login": "lucky"})
		require.ErrorIs(t, resetErr, limiter.ErrIncorrectIdentity)
	})
}
