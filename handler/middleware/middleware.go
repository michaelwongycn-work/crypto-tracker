package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/michaelwongycn/crypto-tracker/lib/auth"
	"github.com/michaelwongycn/crypto-tracker/lib/cache"
)

func Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 {
			http.Error(w, "Malformed token", http.StatusUnauthorized)
			return
		}
		accessToken := tokenParts[1]

		cachedAccessToken := cache.GetCache(accessToken)

		if cachedAccessToken == nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		claims, err := auth.ParseToken(accessToken)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "claims", claims)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
