package handlers

import (
	"encoding/json"
	"github.com/golang-jwt/jwt"
	"github.com/segmentio/ksuid"
	"github.com/soyandrestrujillo/advanced_go_rest_websockets/models"
	"github.com/soyandrestrujillo/advanced_go_rest_websockets/repository"
	"github.com/soyandrestrujillo/advanced_go_rest_websockets/server"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"strings"
	"time"
)

type SignupLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignupResponse struct {
	Id    string `json:"id"`
	Email string `json:"email"`
}

type loginResponse struct {
	Token string `json:"token"`
}

func SignUpHanlder(s server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req = SignupLoginRequest{}
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		id, err := ksuid.NewRandom()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var user = models.User{
			ID:       id.String(),
			Email:    req.Email,
			Password: string(hashedPassword),
		}

		err = repository.InsertUser(r.Context(), &user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(SignupResponse{
			Id:    user.ID,
			Email: user.Email,
		})
	}
}

func LoginHandler(s server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req = SignupLoginRequest{}
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Get user by email
		user, err := repository.GetUserByEmail(r.Context(), req.Email)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Check if user exists
		if user == nil {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
		if err != nil {
			log.Print(err)
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}

		claims := models.AppClaims{
			UserId: user.ID,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: time.Now().Add(2 * time.Hour * 24).Unix(),
			},
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString([]byte(s.Config().JWTSecret))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(loginResponse{
			Token: tokenString,
		})
	}
}

func MyHandler(s server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Don't repeat yourself - start
		tokenString := strings.TrimSpace(r.Header.Get("Authorization"))
		token, err := jwt.ParseWithClaims(tokenString, &models.AppClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(s.Config().JWTSecret), nil
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		// Don't repeat yourself - end

		// Check if token is valid
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
