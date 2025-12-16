package gb

import (
	"time"

	"github.com/rainb0w-clwn/go_auth_limiter/internal/limiter"
)

type TokenBucketGB struct {
	tokenBucketLimiter limiter.ITokenBucketLimitService

	tokenBucketTTL time.Duration
}

func New(tokenBucketLimiter limiter.ITokenBucketLimitService, tokenBucketTTL time.Duration) *TokenBucketGB {
	return &TokenBucketGB{
		tokenBucketLimiter: tokenBucketLimiter,
		tokenBucketTTL:     tokenBucketTTL,
	}
}

func (gb *TokenBucketGB) Sweep() error {
	buckets := gb.tokenBucketLimiter.GetBuckets()
	if len(buckets) == 0 {
		return nil
	}

	bucketsToDelete := make([]string, 0)
	for key, bucket := range buckets {
		if !(*bucket).Full() {
			continue
		}

		lastRefill := (*bucket).GetLastRefill()
		if time.Since(lastRefill) < gb.tokenBucketTTL {
			continue
		}

		bucketsToDelete = append(bucketsToDelete, key)
	}

	if len(bucketsToDelete) == 0 {
		return nil
	}

	for _, bucketKey := range bucketsToDelete {
		err := gb.tokenBucketLimiter.SweepBucket(bucketKey)
		if err != nil {
			return err
		}
	}

	return nil
}

// ITokenBucketGB сервис подчистки устаревших бакетов.
type ITokenBucketGB interface {
	Sweep() error
}
