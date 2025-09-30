package middleware

import (
	"context"
	"github.com/google/uuid"
	"net/http"
)

const requestIdHeaderName = "X-Request-Id"
const RequestIdContextAttributeName = "requestId"

func LoggerContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		xRequestId := r.Header.Get(requestIdHeaderName)
		if xRequestId == "" {
			xRequestId = uuid.NewString()
		}
		ctx = context.WithValue(ctx, RequestIdContextAttributeName, xRequestId)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
