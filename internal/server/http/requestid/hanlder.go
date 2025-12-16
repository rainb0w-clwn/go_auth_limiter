package requestid

import (
	"net/http"

	"github.com/google/uuid"
)

const XRequestIDKey = "X-Request-Id"

func New(next http.Handler) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		requestID := request.Header.Get(XRequestIDKey)
		if requestID == "" {
			requestID = newRequestID()
		}

		writer.Header().Set("X-Request-ID", requestID)

		next.ServeHTTP(writer, request)
	}
}

func newRequestID() string {
	return uuid.NewString()
}
