package handlers

import (
	"context"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"net/http"
)

func AuthMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// TODO: implement auth and user id extraction
			userID := r.Header.Get("Authorization")
			if userID == "" {
				ape.Render(w, problems.Unauthorized())
				return
			}

			r = r.WithContext(context.WithValue(r.Context(), userIDCtxKey, userID))

			next.ServeHTTP(w, r)
		})
	}
}
