package server

import (
	"context"
	"errors"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type Config struct {
	Port        string // Port to run the server on
	JWTSecret   string // Secret to sign JWT tokens with
	DataBaseURL string // URL to connect to the database
}

type Server interface {
	Config() *Config
}

type Broker struct {
	config *Config
	router *mux.Router
}

func (b *Broker) Config() *Config {
	return b.config
}

func NewServer(ctx context.Context, config *Config) (*Broker, error) {
	if config.Port == "" {
		return nil, errors.New("port is required")
	}

	if config.JWTSecret == "" {
		return nil, errors.New("JWTSecret is required")
	}

	if config.DataBaseURL == "" {
		return nil, errors.New("DataBaseURL is required")
	}

	broker := &Broker{
		config: config,
		router: mux.NewRouter(),
	}

	return broker, nil
}

func (b *Broker) Start(binder func(s Server, r *mux.Router)) {
	b.router = mux.NewRouter()
	binder(b, b.router)
	log.Print("Starting server on port " + b.config.Port)
	if err := http.ListenAndServe(b.config.Port, b.router); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
