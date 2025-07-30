package service

import (
	"FocusList/internal/model"
	"FocusList/internal/repository"
	"FocusList/internal/utils"
	"fmt"
	"log"
)

type AuthService struct {
	UserRepo *repository.UserRepository
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
	fmt.Println("Input password:", password)
	user, err := s.UserRepo.GetUserByEmail(email)
	if err != nil {
		log.Println("Error fetching user:", err)
		return "", fmt.Errorf("invalid credentials")
	}
	fmt.Println("Input password:", password)   //TODO: Remove the debug logs
	fmt.Println("Stored hash:", user.Password) //TODO: Remove the debug logs
	if err := utils.CheckPasswordHash(password, user.Password); err != nil {
		return "", fmt.Errorf("invalid credentials")
	}
	return utils.GenerateJWT(user.Email, string(user.Role))
}
