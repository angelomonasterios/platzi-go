package main

import (
	"context"
	"github.com/go/rest-ws/handlers"
	"github.com/go/rest-ws/middleware"
	"github.com/go/rest-ws/server"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
)

func main() {

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file", err)
	}
	PORT := os.Getenv("PORT")
	JWT_SECRET := os.Getenv("JWT_SECRET")
	DATABASE_URL := os.Getenv("DATABASE_URL")

	s, err := server.NewServer(context.Background(), &server.Config{
		JWTSecret:   JWT_SECRET,
		Port:        PORT,
		DataBaseUrl: DATABASE_URL,
	})

	if err != nil {
		log.Fatal(err)
	}

	s.Start(BindRoutes)
}

func BindRoutes(s server.Server, r *mux.Router) {

	api := r.PathPrefix("/api/v1").Subrouter()
	api.Use(middleware.CheckAuthMiddleware(s))
	r.HandleFunc("/", handlers.HomeHandler(s)).Methods(http.MethodGet)
	r.HandleFunc("/signup", handlers.SignUpHandler(s)).Methods(http.MethodPost)
	r.HandleFunc("/login", handlers.LoginHandler(s)).Methods(http.MethodPost)
	api.HandleFunc("/me", handlers.MeHandler(s)).Methods(http.MethodGet)
	api.HandleFunc("/posts", handlers.InsertPostHandler(s)).Methods(http.MethodPost)
	r.HandleFunc("/posts/{id}", handlers.GetPostByIdHandler(s)).Methods(http.MethodGet)
	api.HandleFunc("/posts/{id}", handlers.UpdatePostHandler(s)).Methods(http.MethodPut)
	api.HandleFunc("/posts/{id}", handlers.DeletePostHandler(s)).Methods(http.MethodDelete)
	r.HandleFunc("/posts", handlers.ListPostHandler(s)).Methods(http.MethodGet)
	r.HandleFunc("/ws", s.Hub().HandleWebSocket)
}
