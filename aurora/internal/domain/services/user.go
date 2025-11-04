package services

import (
	"github.com/phankieuphu/go-aws-example/internal/adapters/repositories"

	"github.com/gin-gonic/gin"
)


type UserService struct {
	repository repositories.UserRepository
}

func NewUserServices(repo repositories.UserRepository) *UserService {
	return &UserService{
		repository: repo,
	}
}


func (u *UserService) GetUser(c *gin.Context){
	repo := repositories.NewUserRepository(nil)
	users, err := repo.GetUser(c.Request.Context())
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, users)

}
