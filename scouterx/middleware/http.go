package middleware

import (
	"github.com/zbum/scouter-agent-golang/scouterx/strace"
	"net/http"
)

func HttpTracingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		ctx = strace.StartHttpService(ctx, r)
		defer strace.EndHttpService(ctx, r, r.Response)
		next.ServeHTTP(w, r)
	})
}
