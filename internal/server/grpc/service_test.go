package server_test

import (
	"context"
	"errors"
	"testing"

	mocks "github.com/rainb0w-clwn/go_auth_limiter/internal/interfaces/mocks"
	"github.com/rainb0w-clwn/go_auth_limiter/internal/limiter"
	"github.com/rainb0w-clwn/go_auth_limiter/internal/rule"
	server "github.com/rainb0w-clwn/go_auth_limiter/internal/server/grpc"
	proto "github.com/rainb0w-clwn/go_auth_limiter/proto/limiter"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

//nolint:dupl
func TestService_WhiteListAdd(t *testing.T) {
	ctx := context.Background()
	app := new(mocks.MockApplication)
	logger := new(mocks.MockLogger)
	s := server.NewService(app, logger)

	// успешный вызов
	app.On("WhiteListAdd", "1.2.3.4").Return(nil)
	resp, err := s.WhiteListAdd(ctx, &proto.WhiteListAddRequest{IpNet: "1.2.3.4"})
	require.NoError(t, err)
	require.NotNil(t, resp)
	app.AssertExpectations(t)
	logger.AssertExpectations(t)

	// вызов с ошибкой
	testErr := errors.New("some error")
	app.On("WhiteListAdd", "5.6.7.8").Return(testErr)
	logger.On("Error", mock.Anything).Return()

	resp, err = s.WhiteListAdd(ctx, &proto.WhiteListAddRequest{IpNet: "5.6.7.8"})
	require.Nil(t, resp)
	require.Error(t, err)
	st, _ := status.FromError(err)
	require.Equal(t, codes.Unknown, st.Code())
	require.Contains(t, st.Message(), "some error")

	app.AssertExpectations(t)
	logger.AssertExpectations(t)
}

func TestService_WhiteListDelete(t *testing.T) {
	ctx := context.Background()
	app := new(mocks.MockApplication)
	logger := new(mocks.MockLogger)
	s := server.NewService(app, logger)

	// успешное удаление
	app.On("WhiteListDelete", "1.2.3.4").Return(nil)
	resp, err := s.WhiteListDelete(ctx, &proto.WhiteListDeleteRequest{IpNet: "1.2.3.4"})
	require.NoError(t, err)
	require.NotNil(t, resp)

	// удаление несуществующего правила
	app.On("WhiteListDelete", "5.6.7.8").Return(rule.ErrRuleNotFound)
	logger.On("Error", mock.Anything).Return()

	resp, err = s.WhiteListDelete(ctx, &proto.WhiteListDeleteRequest{IpNet: "5.6.7.8"})
	require.Nil(t, resp)
	require.Error(t, err)
	st, _ := status.FromError(err)
	require.Equal(t, codes.NotFound, st.Code())
	require.Contains(t, st.Message(), rule.ErrRuleNotFound.Error())

	app.AssertExpectations(t)
	logger.AssertExpectations(t)
}

//nolint:dupl
func TestService_BlackListAdd(t *testing.T) {
	ctx := context.Background()
	app := new(mocks.MockApplication)
	logger := new(mocks.MockLogger)
	s := server.NewService(app, logger)

	// успешный вызов
	app.On("BlackListAdd", "1.2.3.4").Return(nil)
	resp, err := s.BlackListAdd(ctx, &proto.BlackListAddRequest{IpNet: "1.2.3.4"})
	require.NoError(t, err)
	require.NotNil(t, resp)
	app.AssertExpectations(t)
	logger.AssertExpectations(t)

	// вызов с ошибкой
	testErr := errors.New("blacklist error")
	app.On("BlackListAdd", "5.6.7.8").Return(testErr)
	logger.On("Error", mock.Anything).Return()

	resp, err = s.BlackListAdd(ctx, &proto.BlackListAddRequest{IpNet: "5.6.7.8"})
	require.Nil(t, resp)
	require.Error(t, err)
	st, _ := status.FromError(err)
	require.Equal(t, codes.Unknown, st.Code())
	require.Contains(t, st.Message(), "blacklist error")

	app.AssertExpectations(t)
	logger.AssertExpectations(t)
}

func TestService_BlackListDelete(t *testing.T) {
	ctx := context.Background()
	app := new(mocks.MockApplication)
	logger := new(mocks.MockLogger)
	s := server.NewService(app, logger)

	// успешное удаление
	app.On("BlackListDelete", "1.2.3.4").Return(nil)
	resp, err := s.BlackListDelete(ctx, &proto.BlackListDeleteRequest{IpNet: "1.2.3.4"})
	require.NoError(t, err)
	require.NotNil(t, resp)

	// удаление несуществующего правила
	app.On("BlackListDelete", "5.6.7.8").Return(rule.ErrRuleNotFound)
	logger.On("Error", mock.Anything).Return()

	resp, err = s.BlackListDelete(ctx, &proto.BlackListDeleteRequest{IpNet: "5.6.7.8"})
	require.Nil(t, resp)
	require.Error(t, err)
	st, _ := status.FromError(err)
	require.Equal(t, codes.NotFound, st.Code())
	require.Contains(t, st.Message(), rule.ErrRuleNotFound.Error())

	app.AssertExpectations(t)
	logger.AssertExpectations(t)
}

func TestService_LimitCheck(t *testing.T) {
	ctx := context.Background()
	app := new(mocks.MockApplication)
	logger := new(mocks.MockLogger)
	s := server.NewService(app, logger)

	// успешная проверка лимита
	app.On("LimitCheck", "1.2.3.4", "user", "pass").Return(true, nil)
	resp, err := s.LimitCheck(ctx, &proto.LimitCheckRequest{Ip: "1.2.3.4", Login: "user", Password: "pass"})
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.True(t, resp.Allowed)
	app.AssertExpectations(t)
	logger.AssertExpectations(t)

	// ошибка неверной идентификации
	app.On("LimitCheck", "1.2.3.4", "user", "wrongpass").Return(false, limiter.ErrIncorrectIdentity)
	logger.On("Error", mock.Anything).Return()

	resp, err = s.LimitCheck(ctx, &proto.LimitCheckRequest{Ip: "1.2.3.4", Login: "user", Password: "wrongpass"})
	require.Nil(t, resp)
	require.Error(t, err)
	st, _ := status.FromError(err)
	require.Equal(t, codes.InvalidArgument, st.Code())
	require.Contains(t, st.Message(), limiter.ErrIncorrectIdentity.Error())

	app.AssertExpectations(t)
	logger.AssertExpectations(t)
}

func TestService_BucketReset(t *testing.T) {
	ctx := context.Background()
	app := new(mocks.MockApplication)
	logger := new(mocks.MockLogger)
	s := server.NewService(app, logger)

	// успешный сброс
	app.On("LimitReset", "1.2.3.4", "user").Return(nil)
	resp, err := s.BucketReset(ctx, &proto.BucketResetRequest{Ip: "1.2.3.4", Login: "user"})
	require.NoError(t, err)
	require.NotNil(t, resp)
	app.AssertExpectations(t)
	logger.AssertExpectations(t)

	// ошибка при сбросе
	testErr := errors.New("reset error")
	app.On("LimitReset", "5.6.7.8", "other").Return(testErr)
	logger.On("Error", mock.Anything).Return()

	resp, err = s.BucketReset(ctx, &proto.BucketResetRequest{Ip: "5.6.7.8", Login: "other"})
	require.Nil(t, resp)
	require.Error(t, err)
	st, _ := status.FromError(err)
	require.Equal(t, codes.Unknown, st.Code())
	require.Contains(t, st.Message(), "reset error")

	app.AssertExpectations(t)
	logger.AssertExpectations(t)
}
