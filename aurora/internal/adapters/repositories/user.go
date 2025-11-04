package repositories

import (
	"context"

	"github.com/phankieuphu/go-aws-example/config"

	"github.com/phankieuphu/go-aws-example/internal/adapters/rds"
	"github.com/phankieuphu/go-aws-example/internal/domain/entities"
	"github.com/phankieuphu/go-aws-example/internal/domain/ports"

	"gorm.io/gorm"
)

type UserRepository struct {
	client *gorm.DB
}

func NewUserRepository(cfg *config.EnvConfig) *UserRepository {
	client := rds.NewRDSClient(cfg)
	return &UserRepository{client: client}
}

var _ ports.UserRepository = (*UserRepository)(nil)

func (u UserRepository) GetUser(ctx context.Context) ([]entities.User, error) {
	var users []entities.User
	if err := u.client.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil

}
