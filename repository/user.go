package repository

import (
	"context"
	"github.com/soyandrestrujillo/advanced_go_rest_websockets/models"
)

// Abstracciones
// Handler - GetUserByIdMongoDB ... ... ... ...

// Concreto (usado)
// Handler - GetUserByI - User ... ... ... ...
// sin importar la base de datos el flujo seguirá siendo el mismo
//           MongoDB
//           MySQL
//           Postgres

type UserRepository interface {
	InsertUser(ctx context.Context, user *models.User) error
	GetUserById(ctx context.Context, id string) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	Close() error
}

var implentation UserRepository

func SetRepository(repository UserRepository) {
	implentation = repository
}

func InsertUser(ctx context.Context, user *models.User) error {
	return implentation.InsertUser(ctx, user)
}

func GetUserById(ctx context.Context, id string) (*models.User, error) {
	return implentation.GetUserById(ctx, id)
}

func GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	return implentation.GetUserByEmail(ctx, email)
}

func Close() error {
	return implentation.Close()
}
