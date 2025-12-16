package token

import (
	"sync"
	"time"

	"github.com/rainb0w-clwn/go_auth_limiter/internal/bucket"
	"github.com/rainb0w-clwn/go_auth_limiter/internal/bucket/refillrate"
)

type Bucket struct {
	sync.RWMutex

	size       int
	refillRate refillrate.RefillRate

	tokensCount int
	lastRefill  time.Time
}

func New(size int, refillRate refillrate.RefillRate) bucket.IBucket {
	return &Bucket{
		size:       size,
		refillRate: refillRate,

		tokensCount: size,
		lastRefill:  time.Now(),
	}
}

func (b *Bucket) Refill() {
	b.Lock()

	const nsInSec = 1e9
	timePassed := time.Since(b.lastRefill)
	tokensToAdd := int64(timePassed) * int64(b.refillRate.GetCount()) / int64(nsInSec*b.refillRate.GetTime().Seconds())

	b.tokensCount = min(b.tokensCount+int(tokensToAdd), b.size)
	if tokensToAdd > 0 {
		b.lastRefill = time.Now()
	}

	b.Unlock()
}

func (b *Bucket) GetTokenCount() int {
	b.RLock()
	defer b.RUnlock()

	return b.tokensCount
}

func (b *Bucket) GetSize() int {
	return b.size
}

func (b *Bucket) GetToken(tokenCount int) {
	b.Lock()

	b.tokensCount -= tokenCount

	b.Unlock()
}

func (b *Bucket) Reset() {
	b.Lock()

	b.tokensCount = b.size
	b.lastRefill = time.Now()

	b.Unlock()
}

func (b *Bucket) GetLastRefill() time.Time {
	return b.lastRefill
}

func (b *Bucket) Full() bool {
	return b.GetTokenCount() == b.GetSize()
}
