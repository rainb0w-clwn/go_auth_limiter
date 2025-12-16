package auth

import (
	"github.com/rainb0w-clwn/go_auth_limiter/internal/limiter"
	"github.com/rainb0w-clwn/go_auth_limiter/internal/limiter/composite"
	"github.com/rainb0w-clwn/go_auth_limiter/internal/rule"
)

// Limiter лимитер для использования в простых формах авторизации с логином и паролем.
type Limiter struct {
	limiter.IService
	ruleService   rule.IService
	bucketLimiter *composite.Limiter
}

func New(
	ruleService rule.IService,
	bucketLimiter *composite.Limiter,
) *Limiter {
	return &Limiter{
		ruleService:   ruleService,
		bucketLimiter: bucketLimiter,
	}
}

func (l *Limiter) SatisfyLimit(identity limiter.UserIdentityDto) (bool, error) {
	validationErr := l.validateIdentity(identity)
	if validationErr != nil {
		return false, validationErr
	}

	inBlackList, blErr := l.ruleService.InBlackList(identity[limiter.IPLimit.String()])
	if inBlackList || blErr != nil {
		return false, blErr
	}

	inWhiteList, wlErr := l.ruleService.InWhiteList(identity[limiter.IPLimit.String()])
	if wlErr != nil {
		return false, wlErr
	}

	if inWhiteList {
		return true, nil
	}

	return l.bucketLimiter.SatisfyLimit(identity)
}

func (l *Limiter) ResetLimit(identity limiter.UserIdentityDto) error {
	return l.bucketLimiter.ResetLimit(identity)
}

func (l *Limiter) SetRequestCost(requestCost int) {
	l.bucketLimiter.SetRequestCost(requestCost)
}

func (l *Limiter) validateIdentity(identity limiter.UserIdentityDto) error {
	if identity[limiter.IPLimit.String()] == "" ||
		identity[limiter.LoginLimit.String()] == "" ||
		identity[limiter.PasswordLimit.String()] == "" {
		// These identity keys are required
		return limiter.ErrIncorrectIdentity
	}

	return nil
}
