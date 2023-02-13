package handlers

import (
	"encoding/json"
	"github.com/go/rest-ws/models"
	"github.com/go/rest-ws/repository"
	"github.com/go/rest-ws/server"
	"github.com/golang-jwt/jwt"
	"net/http"
	"strings"
)

type InsertPostRequest struct {
	PostContent string `json:"post_content"`
}

func InsertPostHandler(s server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := strings.TrimSpace(r.Header.Get("Authorization"))
		var token, err = jwt.ParseWithClaims(tokenString, &models.AppClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(s.Config().JWTSecret), nil
		})

		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		if claims, ok := token.Claims.(*models.AppClaims); ok && token.Valid {
			user, err := repository.GetUserById(r.Context(), claims.UserId)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(user)

		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
