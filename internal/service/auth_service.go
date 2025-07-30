package service

import (
	"FocusList/internal/cache"
	"FocusList/internal/model"
	"FocusList/internal/repository"
	"FocusList/internal/utils"
	"fmt"
	"log"
)

type AuthService struct {
	UserRepo  *repository.UserRepository
	CacheRepo *cache.CacheService
}

func (s *AuthService) Register(user *model.User, plainPassword string) error {
	hash, err := utils.HashPassword(plainPassword)
	if err != nil {
		return err
	}
	user.Password = hash
	user.IsActive = true
	user.Role = "user"
	return s.UserRepo.CreateUser(user)
}

func (s *AuthService) Login(email, password string) (string, error) {
	user, err := s.UserRepo.GetUserByEmail(email)
	if err != nil {
		log.Println("Error fetching user:", err)
		return "", fmt.Errorf("invalid credentials")
	}
	if err := utils.CheckPasswordHash(password, user.Password); err != nil {
		return "", fmt.Errorf("invalid credentials")
	}
	return utils.GenerateJWT(user.Email, string(user.Role))
}
