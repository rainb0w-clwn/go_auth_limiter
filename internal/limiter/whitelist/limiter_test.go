package whitelist_test

import (
	"testing"

	"github.com/rainb0w-clwn/go_auth_limiter/internal/limiter"
	"github.com/rainb0w-clwn/go_auth_limiter/internal/limiter/whitelist"
	"github.com/rainb0w-clwn/go_auth_limiter/internal/rule"
	rulemocks "github.com/rainb0w-clwn/go_auth_limiter/internal/rule/mocks"
	"github.com/stretchr/testify/require"
)

func Test_SatisfyLimit(t *testing.T) {
	rules := rule.Rules{
		rule.Rule{
			ID:       1,
			IP:       "192.168.1.1",
			RuleType: rule.WhiteList,
		},
		rule.Rule{
			ID:       2,
			IP:       "125.130.2.3",
			RuleType: rule.WhiteList,
		},
		rule.Rule{
			ID:       3,
			IP:       "192.168.3.0/24",
			RuleType: rule.WhiteList,
		},
	}
	mockStorage := rulemocks.NewMockIStorage(t)
	mockStorage.EXPECT().GetForType(rule.WhiteList).Return(&rules, nil)

	whiteListLimiter := whitelist.New(
		rule.NewService(mockStorage),
	)

	tests := []struct {
		name     string
		ip       string
		expected bool
	}{
		{
			name:     "found by exact ip match",
			ip:       "192.168.1.1",
			expected: true,
		},
		{
			name:     "found by subnet match",
			ip:       "192.168.3.25",
			expected: true,
		},
		{
			name:     "not found",
			ip:       "222.113.10.3",
			expected: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			satisfies, err := whiteListLimiter.SatisfyLimit(limiter.UserIdentityDto{"ip": test.ip})

			require.NoError(t, err)
			require.Equal(t, test.expected, satisfies)
		})
	}
}

func Test_SatisfyLimit_Error(t *testing.T) {
	t.Run("incorrect identity error", func(t *testing.T) {
		mockStorage := rulemocks.NewMockIStorage(t)
		whiteListLimiter := whitelist.New(
			rule.NewService(mockStorage),
		)

		identity := limiter.UserIdentityDto{"login": "admin"} // white list limiter needs ip
		_, err := whiteListLimiter.SatisfyLimit(identity)
		require.ErrorIs(t, err, limiter.ErrIncorrectIdentity)
	})
}
