package service

import (
	"WebProject/internal/core"
	"context"
)

type UserRepository interface {
	GetAll(ctx context.Context) ([]*core.User, error)
	GetById(ctx context.Context, id string) (*core.User, error)
	Save(ctx context.Context, user *core.User) (*core.User, error)
}

type UserService struct {
	userRepository UserRepository
}

func NewUserService(repository UserRepository) *UserService {
	return &UserService{userRepository: repository}
}

func (service *UserService) GetAll(ctx context.Context) ([]*core.User, error) {
	return service.userRepository.GetAll(ctx)
}

func (service *UserService) GetById(ctx context.Context, id string) (*core.User, error) {
	return service.userRepository.GetById(ctx, id)
}

func (service *UserService) CreateUser(ctx context.Context, user *core.User) (*core.User, error) {
	return service.userRepository.Save(ctx, user)
}
