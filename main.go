package main

import (
	"context"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/soyandrestrujillo/advanced_go_rest_websockets/handlers"
	"github.com/soyandrestrujillo/advanced_go_rest_websockets/server"
	"log"
	"net/http"
	"os"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Printf("Error loading .env file: %v", err)
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
	r.HandleFunc("/signup", handlers.SignUpHanlder(s)).Methods(http.MethodPost)
	r.HandleFunc("/login", handlers.LoginHandler(s)).Methods(http.MethodPost)
}
