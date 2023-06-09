package middleware

import (
	"github.com/zbum/scouter-agent-golang/scouterx/strace"
	"net/http"
)

func FastHttpTracingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		newCtx := strace.StartHttpService(ctx, r)
		newRequest := r.WithContext(newCtx)
		defer strace.EndHttpService(ctx, newRequest, r.Response)
		next.ServeHTTP(w, r)
	})
}
