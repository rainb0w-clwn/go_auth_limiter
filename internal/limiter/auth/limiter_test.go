package auth_test

import (
	"testing"
	"time"

	"github.com/rainb0w-clwn/go_auth_limiter/internal/bucket/refillrate"
	"github.com/rainb0w-clwn/go_auth_limiter/internal/limiter"
	"github.com/rainb0w-clwn/go_auth_limiter/internal/limiter/auth"
	"github.com/rainb0w-clwn/go_auth_limiter/internal/limiter/composite"
	limitermocks "github.com/rainb0w-clwn/go_auth_limiter/internal/limiter/mocks"
	"github.com/rainb0w-clwn/go_auth_limiter/internal/rule"
	rulemocks "github.com/rainb0w-clwn/go_auth_limiter/internal/rule/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestLoginFormLimiter_SatisfyLimit(t *testing.T) {
	limit := 3
	refillRate := refillrate.New(limit, time.Second*1)
	whiteListIP, blackListIP, bothListIP, unknownIP := "192.168.1.1", "192.150.10.3", "10.9.123.12", "5.5.5.5"
	identity := limiter.UserIdentityDto{
		limiter.IPLimit.String():       "unknown",
		limiter.LoginLimit.String():    "lucky",
		limiter.PasswordLimit.String(): "root",
	}

	// Mock RuleService
	ruleStorage := rulemocks.NewMockIStorage(t)
	ruleStorage.EXPECT().GetForType(rule.WhiteList).Return(&rule.Rules{
		rule.Rule{ID: 1, IP: whiteListIP, RuleType: rule.WhiteList},
		rule.Rule{ID: 3, IP: bothListIP, RuleType: rule.WhiteList},
	}, nil).Maybe()
	ruleStorage.EXPECT().GetForType(rule.BlackList).Return(&rule.Rules{
		rule.Rule{ID: 2, IP: blackListIP, RuleType: rule.BlackList},
		rule.Rule{ID: 4, IP: bothListIP, RuleType: rule.BlackList},
	}, nil).Maybe()
	ruleService := rule.NewService(ruleStorage)

	// Mock LimitStorage
	limitStorage := limitermocks.NewMockIStorage(t)
	limitStorage.EXPECT().GetLimitsByTypes(mock.AnythingOfType("[]string")).Return(&limiter.Limits{
		limiter.Limit{LimitType: limiter.IPLimit, Value: limit},
		limiter.Limit{LimitType: limiter.LoginLimit, Value: limit},
		limiter.Limit{LimitType: limiter.PasswordLimit, Value: limit},
	}, nil).Maybe()

	t.Run("ip in white list", func(t *testing.T) {
		loginFormLimiter := auth.New(
			ruleService,
			composite.New(limitStorage, refillRate),
		)
		identity[limiter.IPLimit.String()] = whiteListIP

		satisfies, err := loginFormLimiter.SatisfyLimit(identity)
		require.True(t, satisfies)
		require.NoError(t, err)

		loginFormLimiter.SetRequestCost(limit + 1)
		satisfies, err = loginFormLimiter.SatisfyLimit(identity)
		require.True(t, satisfies)
		require.NoError(t, err)
	})

	t.Run("ip in black list", func(t *testing.T) {
		loginFormLimiter := auth.New(
			ruleService,
			composite.New(limitStorage, refillRate),
		)
		identity[limiter.IPLimit.String()] = blackListIP
		loginFormLimiter.SetRequestCost(limit + 1)

		satisfies, err := loginFormLimiter.SatisfyLimit(identity)
		require.False(t, satisfies)
		require.NoError(t, err)

		loginFormLimiter.SetRequestCost(1)
		satisfies, err = loginFormLimiter.SatisfyLimit(identity)
		require.False(t, satisfies)
		require.NoError(t, err)
	})

	t.Run("ip in both lists", func(t *testing.T) {
		loginFormLimiter := auth.New(
			ruleService,
			composite.New(limitStorage, refillRate),
		)
		identity[limiter.IPLimit.String()] = bothListIP

		satisfies, err := loginFormLimiter.SatisfyLimit(identity)
		require.False(t, satisfies)
		require.NoError(t, err)
	})

	t.Run("no ip in lists", func(t *testing.T) {
		loginFormLimiter := auth.New(
			ruleService,
			composite.New(limitStorage, refillRate),
		)
		identity[limiter.IPLimit.String()] = unknownIP

		satisfies, err := loginFormLimiter.SatisfyLimit(identity)
		require.True(t, satisfies)
		require.NoError(t, err)

		loginFormLimiter.SetRequestCost(limit + 1)
		satisfies, err = loginFormLimiter.SatisfyLimit(identity)
		require.False(t, satisfies)
		require.NoError(t, err)
	})
}

func TestLoginFormLimiter_SatisfyLimit_Error(t *testing.T) {
	ruleService := rule.NewService(rulemocks.NewMockIStorage(t))
	limitStorage := limitermocks.NewMockIStorage(t)
	refillRate := refillrate.New(3, time.Second*1)

	t.Run("incorrect identity", func(t *testing.T) {
		expectedErr := limiter.ErrIncorrectIdentity
		emptyIdentity := limiter.UserIdentityDto{}
		notFullIdentity := limiter.UserIdentityDto{
			limiter.LoginLimit.String():    "lucky",
			limiter.PasswordLimit.String(): "root",
		}
		loginFormLimiter := auth.New(
			ruleService,
			composite.New(limitStorage, refillRate),
		)

		// not full identity #1
		_, err := loginFormLimiter.SatisfyLimit(notFullIdentity)
		require.ErrorIs(t, err, expectedErr)

		// not full identity #2
		delete(notFullIdentity, limiter.LoginLimit.String())
		_, err = loginFormLimiter.SatisfyLimit(notFullIdentity)
		require.ErrorIs(t, err, expectedErr)

		// empty identity
		_, err = loginFormLimiter.SatisfyLimit(emptyIdentity)
		require.ErrorIs(t, err, expectedErr)
	})
}
