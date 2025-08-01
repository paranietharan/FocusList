package service

import (
	"FocusList/internal/cache"
	"FocusList/internal/model"
	"FocusList/internal/repository"
	"FocusList/internal/utils"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"
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
	user.IsActive = false
	user.Role = "user"

	userData, err := json.Marshal(user)
	if err != nil {
		return fmt.Errorf("failed to marshal user data: %w", err)
	}

	userDetailRedisKey := fmt.Sprintf("user:%s", user.Email)
	if err := s.CacheRepo.StoreBytes(userDetailRedisKey, userData, 10*time.Minute); err != nil {
		log.Println("Error storing user data in cache:", err)
	}

	code := utils.GenerateVerificationCode(6)
	userCodeRedisKey := fmt.Sprintf("code:%s", user.Email)
	if err := s.CacheRepo.Set(userCodeRedisKey, code, 10*time.Minute); err != nil {
		log.Println("Error setting verification code in cache:", err)
		return fmt.Errorf("failed to set verification code in cache: %w", err)
	}

	fmt.Printf("Verification code for %s: %s\n", user.Email, code)
	return nil
}

func (s *AuthService) VerifyEmail(email, code string) error {
	userCodeRedisKey := fmt.Sprintf("code:%s", email)
	storedCode, err := s.CacheRepo.Get(userCodeRedisKey)
	if err != nil {
		log.Println("Error getting verification code from cache:", err)
		return fmt.Errorf("failed to get verification code from cache: %w", err)
	}

	// TODO: Remove the debug logs
	fmt.Printf("Stored code for %s: %s\n", email, storedCode)
	fmt.Printf("Provided code for %s: %s\n", email, code)

	if storedCode == "" {
		log.Println("No verification code found for email in cache:", email)
		return fmt.Errorf("no verification code found for email: %s", email)
	}

	sc, err1 := strconv.Atoi(storedCode)
	if err1 != nil {
		log.Println("Error converting stored code to integer:", err1)
		return fmt.Errorf("invalid verification code format")
	}
	c, err2 := strconv.Atoi(code)
	if err2 != nil {
		log.Println("Error converting provided code to integer:", err2)
		return fmt.Errorf("invalid verification code format")
	}

	if sc != c {
		log.Println("Verification code mismatch")
		return fmt.Errorf("verification code does not match")
	}

	userDetailRedisKey := fmt.Sprintf("user:%s", email)
	userData, err := s.CacheRepo.GetBytes(userDetailRedisKey)
	if err != nil {
		log.Println("Error getting user data from cache:", err)
		return fmt.Errorf("failed to get user data from cache: %w", err)
	}

	var user model.User
	if err := json.Unmarshal(userData, &user); err != nil {
		log.Println("Error unmarshalling user data:", err)
		return fmt.Errorf("failed to unmarshal user data: %w", err)
	}

	user.IsActive = true
	if err := s.UserRepo.CreateUser(&user); err != nil {
		log.Println("Error updating user:", err)
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
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
