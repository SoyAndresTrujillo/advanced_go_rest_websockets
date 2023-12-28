package main

import (
	"context"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/soyandrestrujillo/rest-ws/handlers"
	"github.com/soyandrestrujillo/rest-ws/server"
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
	JWTSecret := os.Getenv("JWT_SECRET")
	DatabaseUrl := os.Getenv("DATABASE_URL")

	s, err := server.NewServer(context.Background(), &server.Config{
		JWTSecret:   JWTSecret,
		Port:        PORT,
		DataBaseURL: DatabaseUrl,
	})

	if err != nil {
		log.Fatal("NewServer: ", err)
	}

	s.Start(BindRoutes)
}

func BindRoutes(s server.Server, r *mux.Router) {
	r.HandleFunc("/", handlers.HomeHandler(s)).Methods(http.MethodGet)
}
