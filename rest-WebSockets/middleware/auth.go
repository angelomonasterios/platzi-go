package middleware

import (
	"github.com/go/rest-ws/models"
	"github.com/go/rest-ws/server"
	"github.com/golang-jwt/jwt"
	"net/http"
	"strings"
)

var (
	NO_AUTH_NEEDED = []string{
		"login",
		"signup",
	}
)

func shouldCheckTocken(route string) bool {
	for _, p := range NO_AUTH_NEEDED {
		if strings.Contains(route, p) {
			return false
		}
	}
	return true

}

func CheckAuthMiddleware(s server.Server) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !shouldCheckTocken(r.URL.Path) {
				next.ServeHTTP(w, r)
				return
			}
			tokenStirng := strings.TrimSpace(r.Header.Get("Authorization"))
			var _, err = jwt.ParseWithClaims(tokenStirng, &models.AppClaims{}, func(token *jwt.Token) (interface{}, error) {
				return []byte(s.Config().JWTSecret), nil
			})

			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
