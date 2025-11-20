package middleware

import (
	"context"
	"net/http"
	"github.com/google/uuid"
)

type ctxKey string

const requestIDKey ctxKey = "requestID"

func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqID := r.Header.Get("X-Request-ID")
		if reqID == "" {
			reqID = uuid.New().String()
		}

		ctx := context.WithValue(r.Context(), requestIDKey, reqID)

		w.Header().Set("X-Request-ID", reqID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetRequestID(r *http.Request) string {
	if val, ok := r.Context().Value(requestIDKey).(string); ok {
		return val
	}
	return ""
}
