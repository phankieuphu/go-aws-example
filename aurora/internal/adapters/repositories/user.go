package repositories

import (
	"context"
	"ex-aurora/config"
	"ex-aurora/internal/adapters/rds"
	"ex-aurora/internal/domain/entities"
	"ex-aurora/internal/domain/ports"

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
