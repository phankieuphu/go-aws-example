package ports

import (
	"context"
	"ex-aurora/internal/domain/entities"
)

type UserRepository interface {
	GetUser(ctx context.Context) ([]entities.User, error)
}
