package ports

import (
	"context"

	"github.com/phankieuphu/go-aws-example/internal/domain/entities"
)

type UserRepository interface {
	GetUser(ctx context.Context) ([]entities.User, error)
}
