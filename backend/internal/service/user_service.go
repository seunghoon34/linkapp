package service

import (
	"errors"

	"github.com/seunghoon34/linkapp/backend/internal/model"
	"github.com/seunghoon34/linkapp/backend/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) CreateUser(username, email, password string) (*model.User, error) {
	// Check if user already exists
	existingUser, _ := s.repo.GetByEmail(email)
	if existingUser != nil {
		return nil, errors.New("user with this email already exists")
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		Username: username,
		Email:    email,
		Password: string(hashedPassword),
	}

	err = s.repo.Create(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) GetUserByID(id string) (*model.User, error) {
	return s.repo.GetByID(id)
}

func (s *UserService) GetUserByEmail(email string) (*model.User, error) {
	return s.repo.GetByEmail(email)
}

func (s *UserService) UpdateUser(id string, username, email string) error {
	user, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}

	user.Username = username
	user.Email = email

	return s.repo.Update(user)
}

func (s *UserService) AuthenticateUser(email, password string) (*model.User, error) {
	user, err := s.repo.GetByEmail(email)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	return user, nil
}
