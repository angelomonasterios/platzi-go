package server

import (
	"context"
	"errors"
	"github.com/go/rest-ws/database"
	"github.com/go/rest-ws/repository"
	"github.com/go/rest-ws/webSocket"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"log"
	"net/http"
)

type Config struct {
	Port        string
	JWTSecret   string
	DataBaseUrl string
}

type Server interface {
	Config() *Config
	Hub() *webSocket.Hub
}

type Broker struct {
	config *Config
	router *mux.Router
	hub    *webSocket.Hub
}

func (b *Broker) Config() *Config {
	return b.config
}

func (b *Broker) Hub() *webSocket.Hub {
	return b.hub
}

func NewServer(ctx context.Context, config *Config) (*Broker, error) {
	if config.Port == "" {
		return nil, errors.New("port is required")
	}
	if config.JWTSecret == "" {
		return nil, errors.New("JWTSecret is required")
	}

	if config.DataBaseUrl == "" {
		return nil, errors.New("DataBaseUrl is required")
	}

	broker := &Broker{
		config: config,
		router: mux.NewRouter(),
		hub:    webSocket.NewHub(),
	}
	return broker, nil
}

func (b *Broker) Start(binder func(s Server, r *mux.Router)) {
	b.router = mux.NewRouter()
	binder(b, b.router)
	hanlder := cors.Default().Handler(b.router)

	repo, err := database.NewPostgresRepository(b.config.DataBaseUrl)

	if err != nil {
		log.Fatal(err)
	}
	go b.hub.Run()
	repository.SetRepository(repo)

	log.Println("Starting server on port: ", b.Config().Port)

	if err := http.ListenAndServe(b.config.Port, hanlder); err != nil {
		log.Fatal("ListendAndServer: ", err)
	}

}
