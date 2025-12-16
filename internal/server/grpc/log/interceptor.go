package log

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"github.com/rainb0w-clwn/go_auth_limiter/internal/interfaces"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

const XRequestIDKey = "x-request-id"

const unknown = "UNKNOWN"

func New(logger interfaces.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		start := time.Now()
		resp, err := handler(ctx, req)
		end := time.Since(start)
		headers, ok := metadata.FromIncomingContext(ctx)

		ip := unknown
		requestID := unknown
		peerInfo, peerOk := peer.FromContext(ctx)
		if peerOk {
			ip = peerInfo.Addr.String()
		}
		if ok { //nolint:nestif
			xRequestID := headers.Get(XRequestIDKey)
			if len(xRequestID) > 0 && xRequestID[0] != "" {
				requestID = xRequestID[0]
			}
			if !peerOk {
				xForwardFor := headers.Get("x-forwarded-for")
				if len(xForwardFor) > 0 && xForwardFor[0] != "" {
					ips := strings.Split(xForwardFor[0], ",")
					if len(ips) > 0 {
						ip = ips[0]
					}
				}
			}
		}

		userAgent := unknown
		if ok {
			userAgent = headers.Get("user-agent")[0]
		}

		statusCode := codes.Unknown
		if st, ok := status.FromError(err); ok {
			statusCode = st.Code()
		}

		logJSON, marshalErr := json.Marshal(
			struct {
				RequestID string
				IP        string
				Datetime  string
				Method    string
				Status    string
				Time      string
				UserAgent string
			}{
				RequestID: requestID,
				IP:        ip,
				Datetime:  time.Now().Format(time.RFC822),
				Method:    info.FullMethod,
				Status:    strconv.Itoa(int(statusCode)),
				Time:      end.String(),
				UserAgent: userAgent,
			},
		)
		if marshalErr != nil {
			logger.Error(marshalErr.Error())
		}

		logger.Info(string(logJSON))

		return resp, err
	}
}
