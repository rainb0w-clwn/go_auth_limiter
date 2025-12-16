package tokenbucket

import (
	"sync"

	"github.com/rainb0w-clwn/go_auth_limiter/internal/bucket"
	"github.com/rainb0w-clwn/go_auth_limiter/internal/bucket/refillrate"
	"github.com/rainb0w-clwn/go_auth_limiter/internal/bucket/token"
	"github.com/rainb0w-clwn/go_auth_limiter/internal/limiter"
)

const DefaultRequestCost = 1

// TokenBucketLimiter позволяет задать rate limit для запросов с использованием алгоритма Bucket.
//
// Для идентификации клиента запроса используется обобщенный объект UserIdentityDto.
// Ключ bucketKey используется для поиска идентификатора клиента в UserIdentityDto.
//
// Позволяет проверять возможность выполнения очередного запроса и получать количество доступных.
type Limiter struct {
	sync.Mutex

	buckets    map[string]*bucket.IBucket
	bucketSize int

	// Скорость пополнения токенов корзины.
	bucketRefillRate refillrate.RefillRate

	// Количество токенов, которое тратится на один запрос при вызове SatisfyLimit.
	// По умолчанию - 1.
	requestCost int

	// Ключ корзины. По данному ключу происходит поиск идентификационных данных методом SatisfyLimit и ResetLimit.
	bucketKey string
}

func New(bucketKey string, bucketSize int, refillRate refillrate.RefillRate) limiter.ITokenBucketLimitService {
	return &Limiter{
		buckets:     make(map[string]*bucket.IBucket),
		requestCost: DefaultRequestCost,

		bucketKey:        bucketKey,
		bucketSize:       bucketSize,
		bucketRefillRate: refillRate,
	}
}

// SatisfyLimit проверяет возможность выполнения запроса для identity c учетом текущей стоимости запроса.
// Происходит забор токенов из корзины для identity c учетом текущей стоимости запроса.
// Выполняется пополнение корзины для identity с учетом заданной скорости пополнения.
func (l *Limiter) SatisfyLimit(identity limiter.UserIdentityDto) (bool, error) {
	identityValue, found := identity[l.bucketKey]
	if !found {
		return false, limiter.ErrIncorrectIdentity
	}

	l.Lock()
	defer l.Unlock()
	b := l.initBucket(identityValue)

	(*b).Refill()

	if (*b).GetTokenCount() > 0 && l.requestCost <= (*b).GetTokenCount() {
		(*b).GetToken(l.requestCost)

		return true, nil
	}

	return false, nil
}

func (l *Limiter) ResetLimit(identity limiter.UserIdentityDto) error {
	identityValue, foundIdentityValue := identity[l.bucketKey]
	if !foundIdentityValue {
		return limiter.ErrIncorrectIdentity
	}

	l.Lock()
	defer l.Unlock()
	b := l.findBucket(identityValue)
	if b == nil {
		return nil
	}

	(*b).Reset()

	return nil
}

func (l *Limiter) SweepBucket(bucketKey string) error {
	l.Lock()

	delete(l.GetBuckets(), bucketKey)

	l.Unlock()

	return nil
}

func (l *Limiter) SetRequestCost(requestCost int) {
	l.requestCost = requestCost
}

// GetRequestsAllowed возвращает количество возможных запросов для identity с учетом текущей стоимости запроса.
func (l *Limiter) GetRequestsAllowed(identity limiter.UserIdentityDto) (int, error) {
	identityValue, found := identity[l.bucketKey]
	if !found {
		return 0, limiter.ErrIncorrectIdentity
	}

	l.Lock()
	defer l.Unlock()

	b := l.initBucket(identityValue)

	(*b).Refill()

	return (*b).GetTokenCount() / l.requestCost, nil
}

func (l *Limiter) GetBuckets() map[string]*bucket.IBucket {
	return l.buckets
}

func (l *Limiter) initBucket(identityValue string) *bucket.IBucket {
	b := l.findBucket(identityValue)
	if b == nil {
		b = l.createBucket(identityValue)
	}

	return b
}

func (l *Limiter) findBucket(identityValue string) *bucket.IBucket {
	b, found := l.buckets[identityValue]
	if !found {
		return nil
	}

	return b
}

func (l *Limiter) createBucket(identityValue string) *bucket.IBucket {
	newBucket := token.New(l.bucketSize, l.bucketRefillRate)
	l.buckets[identityValue] = &newBucket

	return l.buckets[identityValue]
}
