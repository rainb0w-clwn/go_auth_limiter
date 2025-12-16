package composite

import (
	"errors"
	"math"
	"strings"

	"github.com/rainb0w-clwn/go_auth_limiter/internal/bucket"
	"github.com/rainb0w-clwn/go_auth_limiter/internal/bucket/refillrate"
	"github.com/rainb0w-clwn/go_auth_limiter/internal/limiter"
	"github.com/rainb0w-clwn/go_auth_limiter/internal/limiter/tokenbucket"
)

var ErrNoLimitsFound = errors.New("not found any limits for given identity")

// Limiter лимитер с использованием нескольких bucket'ов
// Набор bucket'ов определяется на основе входных данных в UserIdentityDto (ключей).
// Объединение по логике И: для удовлетворения лимиту необходимо "пройти" все bucket'ы.
type Limiter struct {
	limitStorage limiter.IStorage

	limiters map[string]limiter.ITokenBucketLimitService

	refillRate  refillrate.RefillRate
	requestCost int
}

func New(limitStorage limiter.IStorage, refillRate refillrate.RefillRate) *Limiter {
	return &Limiter{
		limitStorage: limitStorage,
		refillRate:   refillRate,
		requestCost:  tokenbucket.DefaultRequestCost,
	}
}

func (o *Limiter) SatisfyLimit(identity limiter.UserIdentityDto) (bool, error) {
	identityKeys := o.getIdentityKeys(identity)
	if len(identity) == 0 {
		return false, limiter.ErrIncorrectIdentity
	}

	if len(o.limiters) == 0 {
		limitersInitErr := o.init(identityKeys)
		if limitersInitErr != nil {
			return false, limitersInitErr
		}
	}

	for key := range identity {
		l, found := o.limiters[key]
		if !found {
			return false, limiter.ErrIncorrectIdentity
		}

		satisfies, checkErr := l.SatisfyLimit(identity)
		if checkErr != nil {
			return false, checkErr
		}

		if !satisfies {
			return false, nil // not satisfies if fails at least one limiter
		}
	}

	return true, nil // satisfies if pass all limiter
}

func (o *Limiter) ResetLimit(identity limiter.UserIdentityDto) error {
	if len(identity) == 0 {
		return limiter.ErrIncorrectIdentity
	}

	if len(o.limiters) == 0 {
		return ErrNoLimitsFound
	}

	for key := range identity {
		l, found := o.limiters[key]
		if !found {
			return limiter.ErrIncorrectIdentity
		}

		err := l.ResetLimit(identity)
		if err != nil {
			return err
		}
	}

	return nil
}

func (o *Limiter) SweepBucket(compositeKey string) error {
	limiterKey, bucketKey, foundSep := strings.Cut(compositeKey, "_")
	if !foundSep {
		return limiter.ErrIncorrectBucketKey
	}

	l, foundLimiter := o.limiters[limiterKey]
	if !foundLimiter {
		return limiter.ErrIncorrectBucketKey
	}

	return l.SweepBucket(bucketKey)
}

func (o *Limiter) SetRequestCost(requestCost int) {
	o.requestCost = requestCost

	if len(o.limiters) == 0 {
		return
	}

	for _, l := range o.limiters {
		l.SetRequestCost(requestCost)
	}
}

// GetRequestsAllowed возращает минимум из остатков всех лимитеров.
func (o *Limiter) GetRequestsAllowed(identity limiter.UserIdentityDto) (int, error) {
	identityKeys := o.getIdentityKeys(identity)
	if len(identity) == 0 {
		return 0, limiter.ErrIncorrectIdentity
	}

	if len(o.limiters) == 0 {
		limitersInitErr := o.init(identityKeys)
		if limitersInitErr != nil {
			return 0, limitersInitErr
		}
	}

	minAllowed := math.MaxInt
	for key := range identity {
		l, found := o.limiters[key]
		if !found {
			return 0, limiter.ErrIncorrectIdentity
		}

		limiterAllowed, checkErr := l.GetRequestsAllowed(identity)
		if checkErr != nil {
			return 0, checkErr
		}

		minAllowed = min(minAllowed, limiterAllowed)
	}

	return minAllowed, nil
}

func (o *Limiter) GetBuckets() map[string]*bucket.IBucket {
	buckets := make(map[string]*bucket.IBucket)

	for limiterKey, l := range o.limiters {
		limiterBuckets := l.GetBuckets()
		for bucketKey, b := range limiterBuckets {
			buckets[limiterKey+"_"+bucketKey] = b
		}
	}

	return buckets
}

func (o *Limiter) init(identityKeys []string) error {
	limits, getLimitsErr := o.limitStorage.GetLimitsByTypes(identityKeys)
	if getLimitsErr != nil || len(*limits) == 0 {
		return ErrNoLimitsFound
	}

	o.limiters = make(map[string]limiter.ITokenBucketLimitService, len(*limits))
	for _, limit := range *limits {
		key := limit.LimitType.String()
		l := tokenbucket.New(key, limit.Value, o.refillRate)
		l.SetRequestCost(o.requestCost)

		o.limiters[key] = l
	}

	return nil
}

func (o *Limiter) getIdentityKeys(identity limiter.UserIdentityDto) []string {
	keys := make([]string, 0, len(identity))

	for key := range identity {
		keys = append(keys, key)
	}

	return keys
}
