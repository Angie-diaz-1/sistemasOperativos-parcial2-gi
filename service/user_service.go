package service

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"usersProject/models"
	"usersProject/repository"
)

type UserServiceInterface interface {
	CreateUser(ctx context.Context, user *models.User) (*mongo.InsertOneResult, error)
	GetAllUsers(ctx context.Context) ([]models.User, error)
	GetUserByID(ctx context.Context, id string) (*models.User, error)
	UpdateUser(ctx context.Context, id string, user *models.User) (*mongo.UpdateResult, error)
	DeleteUser(ctx context.Context, id string) (*mongo.DeleteResult, error)
}

type UserService struct {
	repo repository.UserRepositoryInterface
}

func NewUserService(repo repository.UserRepositoryInterface) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) CreateUser(ctx context.Context, user *models.User) (*mongo.InsertOneResult, error) {
	return s.repo.Create(ctx, user)
}

func (s *UserService) GetAllUsers(ctx context.Context) ([]models.User, error) {
	return s.repo.GetAll(ctx)
}

func (s *UserService) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *UserService) UpdateUser(ctx context.Context, id string, user *models.User) (*mongo.UpdateResult, error) {
	return s.repo.Update(ctx, id, user)
}

func (s *UserService) DeleteUser(ctx context.Context, id string) (*mongo.DeleteResult, error) {
	return s.repo.Delete(ctx, id)
}
