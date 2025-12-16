package limiter

import (
	"errors"

	"github.com/rainb0w-clwn/go_auth_limiter/internal/bucket"
)

const (
	LoginLimit    Type = "login"
	PasswordLimit Type = "password"
	IPLimit       Type = "ip"
)

var (
	// ErrIncorrectIdentity Ошибка на случай некорректного входного аргумента identity.
	ErrIncorrectIdentity  = errors.New("not found appropriate key in user identity")
	ErrNotSupported       = errors.New("operation not supported")
	ErrIncorrectBucketKey = errors.New("incorrect bucket key")
)

type Type string

func (t Type) String() string {
	return string(t)
}

type Limits []Limit

type Limit struct {
	LimitType   Type
	Value       int
	Description string
}

// IStorage хранилище лимитов (правил) rate limit'инга запросов.
type IStorage interface {
	GetLimits() (*Limits, error)
	GetLimitsByTypes([]string) (*Limits, error)
}

// IService основной сервис проверки запроса на rate limit.
type IService interface {
	SatisfyLimit(UserIdentityDto) (bool, error)
	ResetLimit(UserIdentityDto) error
}

// UserIdentityDto тип для идентификации клиента, запрос которого подвергается rate limit'ингу.
// Может содержать один или несколько пар ключ-значение. Лимитеры сами решают, с какими ключами работать.
type UserIdentityDto map[string]string

// ITokenBucketLimitService интерфейс лимитеров на основе Bucket.
type ITokenBucketLimitService interface {
	IService

	SetRequestCost(int)
	GetRequestsAllowed(UserIdentityDto) (int, error)
	GetBuckets() map[string]*bucket.IBucket
	SweepBucket(string) error
}
