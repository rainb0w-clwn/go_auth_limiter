package rule_test

import (
	"errors"
	"testing"

	"github.com/rainb0w-clwn/go_auth_limiter/internal/rule"
	rulemocks "github.com/rainb0w-clwn/go_auth_limiter/internal/rule/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func newService(t *testing.T) (*rule.Service, *rulemocks.MockIStorage) {
	t.Helper()

	storage := rulemocks.NewMockIStorage(t)
	service := rule.NewService(storage)

	return service, storage
}

func TestService_InList(t *testing.T) {
	errDB := errors.New("db error")
	tests := []struct {
		name        string
		listType    rule.Type
		inputIP     string
		rules       rule.Rules
		storageErr  error
		expected    bool
		expectedErr error
	}{
		{
			name:     "direct ip match",
			listType: rule.WhiteList,
			inputIP:  "192.168.1.10",
			rules: rule.Rules{
				{ID: 1, IP: "192.168.1.10", RuleType: rule.WhiteList},
			},
			expected: true,
		},
		{
			name:     "cidr match",
			listType: rule.BlackList,
			inputIP:  "10.0.0.42",
			rules: rule.Rules{
				{ID: 1, IP: "10.0.0.0/24", RuleType: rule.BlackList},
			},
			expected: true,
		},
		{
			name:     "not found",
			listType: rule.WhiteList,
			inputIP:  "192.168.1.99",
			rules: rule.Rules{
				{ID: 1, IP: "192.168.1.1", RuleType: rule.WhiteList},
			},
			expected: false,
		},
		{
			name:        "invalid input ip",
			listType:    rule.WhiteList,
			inputIP:     "invalid-ip",
			expected:    false,
			expectedErr: rule.ErrInvalidInputIP,
		},
		{
			name:        "storage error",
			listType:    rule.BlackList,
			inputIP:     "127.0.0.1",
			storageErr:  errDB,
			expected:    false,
			expectedErr: errDB,
		},
		{
			name:     "invalid cidr in storage is ignored",
			listType: rule.WhiteList,
			inputIP:  "10.0.0.1",
			rules: rule.Rules{
				{ID: 1, IP: "not-a-cidr", RuleType: rule.WhiteList},
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, storage := newService(t)

			// invalid IP — storage НЕ должен вызываться
			if errors.Is(tt.expectedErr, rule.ErrInvalidInputIP) {
				var result bool
				var err error

				if tt.listType == rule.WhiteList {
					result, err = service.InWhiteList(tt.inputIP)
				} else {
					result, err = service.InBlackList(tt.inputIP)
				}

				require.False(t, result)
				require.ErrorIs(t, err, rule.ErrInvalidInputIP)

				storage.AssertNotCalled(t, "GetForType", mock.Anything)
				return
			}

			storage.
				On("GetForType", tt.listType).
				Return(func() *rule.Rules {
					if tt.storageErr != nil {
						return nil
					}
					return &tt.rules
				}(), tt.storageErr).
				Once()

			var result bool
			var err error

			if tt.listType == rule.WhiteList {
				result, err = service.InWhiteList(tt.inputIP)
			} else {
				result, err = service.InBlackList(tt.inputIP)
			}

			if tt.expectedErr != nil {
				require.ErrorIs(t, err, tt.expectedErr)
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, tt.expected, result)

			storage.AssertExpectations(t)
		})
	}
}

func TestService_ListAdd(t *testing.T) {
	tests := []struct {
		name     string
		listType rule.Type
	}{
		{"white list add", rule.WhiteList},
		{"black list add", rule.BlackList},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, storage := newService(t)

			storage.
				On("Create", rule.Rule{
					IP:       "127.0.0.1",
					RuleType: tt.listType,
				}).
				Return(1, nil).
				Once()

			var err error
			if tt.listType == rule.WhiteList {
				err = service.WhiteListAdd("127.0.0.1")
			} else {
				err = service.BlackListAdd("127.0.0.1")
			}

			require.NoError(t, err)
			storage.AssertExpectations(t)
		})
	}
}

func TestService_ListDelete(t *testing.T) {
	someError := errors.New("some error")
	tests := []struct {
		name        string
		listType    rule.Type
		rules       rule.Rules
		findErr     error
		deleteErr   error
		expectedErr error
	}{
		{
			name:     "success delete all",
			listType: rule.WhiteList,
			rules: rule.Rules{
				{ID: 1},
				{ID: 2},
			},
		},
		{
			name:        "rule not found",
			listType:    rule.BlackList,
			rules:       rule.Rules{},
			expectedErr: rule.ErrRuleNotFound,
		},
		{
			name:        "find error",
			listType:    rule.WhiteList,
			findErr:     someError,
			expectedErr: someError,
		},
		{
			name:        "delete error",
			listType:    rule.BlackList,
			rules:       rule.Rules{{ID: 1}},
			deleteErr:   someError,
			expectedErr: someError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, storage := newService(t)

			storage.
				On("Find", "127.0.0.1", tt.listType).
				Return(func() *rule.Rules {
					if tt.findErr != nil {
						return nil
					}
					return &tt.rules
				}(), tt.findErr).
				Once()

			if tt.findErr == nil {
				for _, r := range tt.rules {
					storage.
						On("Delete", r.ID).
						Return(tt.deleteErr).
						Once()

					if tt.deleteErr != nil {
						break
					}
				}
			}

			var err error
			if tt.listType == rule.WhiteList {
				err = service.WhiteListDelete("127.0.0.1")
			} else {
				err = service.BlackListDelete("127.0.0.1")
			}

			if tt.expectedErr != nil {
				require.ErrorIs(t, err, tt.expectedErr)
			} else {
				require.NoError(t, err)
			}

			storage.AssertExpectations(t)
		})
	}
}
