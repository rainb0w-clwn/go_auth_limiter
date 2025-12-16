package whitelist

import (
	"github.com/rainb0w-clwn/go_auth_limiter/internal/limiter"
	"github.com/rainb0w-clwn/go_auth_limiter/internal/rule"
)

const IdentityKey = "ip"

// Limiter лимитер на основе принципа белого списка.
// Если клиент есть в белом списке, то считается, что он всегда удовлетворяет лимитам.
type Limiter struct {
	ruleService rule.IService
}

func New(service *rule.Service) limiter.IService {
	return &Limiter{
		ruleService: service,
	}
}

func (l Limiter) SatisfyLimit(identity limiter.UserIdentityDto) (bool, error) {
	ip, found := identity[IdentityKey]
	if !found {
		return false, limiter.ErrIncorrectIdentity
	}

	inWhiteList, err := l.ruleService.InWhiteList(ip)
	if err != nil {
		return false, err
	}

	return inWhiteList, nil
}

func (l Limiter) ResetLimit(_ limiter.UserIdentityDto) error {
	return limiter.ErrNotSupported
}
