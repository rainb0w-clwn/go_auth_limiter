package log

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/rainb0w-clwn/go_auth_limiter/internal/interfaces"
)

const XRequestIDKey = "X-Request-Id"

func New(logger appinterfaces.Logger, next http.Handler) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if request.URL.Path == "/health" {
			return
		}
		lrw := &loggingResponseWriter{writer, http.StatusOK}

		start := time.Now()
		next.ServeHTTP(writer, request)
		end := time.Since(start)

		logJSON, err := json.Marshal(
			struct {
				RequestID string
				IP        string
				Datetime  string
				Method    string
				Path      string
				HTTP      string
				Status    string
				Time      string
				UserAgent string
			}{
				RequestID: request.Header.Get(XRequestIDKey),
				IP:        request.RemoteAddr,
				Datetime:  time.Now().Format(time.RFC822),
				Method:    request.Method,
				Path:      request.URL.Path,
				HTTP:      request.Proto,
				Status:    strconv.Itoa(lrw.StatusCode),
				Time:      end.String(),
				UserAgent: request.UserAgent(),
			},
		)
		if err != nil {
			logger.Error(err.Error())
		}

		logger.Info(string(logJSON))
	}
}

type loggingResponseWriter struct {
	http.ResponseWriter
	StatusCode int
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.StatusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}
