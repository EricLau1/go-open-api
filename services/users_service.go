package services

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"go-open-api/models"
	"go-open-api/repository"
	"go-open-api/security/passwords"
	"go-open-api/types"
	"time"
)

type UsersService interface {
	CreateUser(ctx context.Context, in *types.CreateUserInput) (*types.User, error)
	UpdateEmail(ctx context.Context, in *types.UpdateEmailInput) (*types.User, error)
	UpdatePassword(ctx context.Context, in *types.UpdatePasswordInput) (*types.User, error)
	GetUser(ctx context.Context, in *types.GetUserInput) (*types.User, error)
	GetUsers(ctx context.Context) ([]*types.User, error)
	DeleteUser(ctx context.Context, in *types.DeleteUserInput) error
}

type usersService struct {
	usersRepository repository.UsersRepository
}

func NewUsersService(usersRepository repository.UsersRepository) UsersService {
	return &usersService{usersRepository: usersRepository}
}

func (s *usersService) CreateUser(ctx context.Context, in *types.CreateUserInput) (*types.User, error) {

	user := &models.User{
		ID:        uuid.NewString(),
		Email:     in.Email,
		Password:  passwords.New(in.Password),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := s.usersRepository.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	return user.ToResponse(), nil
}

func (s *usersService) UpdateEmail(ctx context.Context, in *types.UpdateEmailInput) (*types.User, error) {
	user, err := s.usersRepository.Get(ctx, in.ID)
	if err != nil {
		return nil, err
	}

	if !passwords.Verify(user.Password, in.Password) {
		return nil, fmt.Errorf("cannot update: %s", user.Email)
	}

	user.Email = in.Email
	user.UpdatedAt = time.Now()

	err = s.usersRepository.Update(ctx, user)
	if err != nil {
		return nil, err
	}

	return user.ToResponse(), nil
}

func (s *usersService) UpdatePassword(ctx context.Context, in *types.UpdatePasswordInput) (*types.User, error) {
	user, err := s.usersRepository.Get(ctx, in.ID)
	if err != nil {
		return nil, err
	}

	if !passwords.Verify(user.Password, in.CurrentPassword) {
		return nil, fmt.Errorf("cannot update: %s", user.Email)
	}

	user.Password = passwords.New(in.NewPassword)
	user.UpdatedAt = time.Now()

	err = s.usersRepository.Update(ctx, user)
	if err != nil {
		return nil, err
	}

	return user.ToResponse(), nil
}

func (s *usersService) GetUser(ctx context.Context, in *types.GetUserInput) (*types.User, error) {
	user, err := s.usersRepository.Get(ctx, in.ID)
	if err != nil {
		return nil, err
	}
	return user.ToResponse(), nil
}

func (s *usersService) GetUsers(ctx context.Context) ([]*types.User, error) {
	users, err := s.usersRepository.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	res := make([]*types.User, 0, len(users))

	for index := range users {
		res = append(res, users[index].ToResponse())
	}

	return res, nil
}

func (s *usersService) DeleteUser(ctx context.Context, in *types.DeleteUserInput) error {
	return s.usersRepository.Delete(ctx, in.ID)
}
