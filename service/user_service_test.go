package service

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"usersProject/models"
)

type MockUserRepository struct {
	mock.Mock
}

// Puedes agregar más métodos si los vas a testear:
func (m *MockUserRepository) GetAll(ctx context.Context) ([]models.User, error) {
	args := m.Called(ctx)
	return args.Get(0).([]models.User), args.Error(1)
}

func (m *MockUserRepository) GetByID(ctx context.Context, id string) (*models.User, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*models.User), args.Error(1)
}

func TestGetAllUsers(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	expectedUsers := []models.User{
		{Name: "Juan", LastName: "Henao", Email: "juan@example.com", Password: "123"},
		{Name: "Angie", LastName: "Diaz", Email: "angie@example.com", Password: "456"},
	}

	mockRepo.On("GetAll", mock.Anything).Return(expectedUsers, nil)

	result, err := service.GetAllUsers(context.TODO())

	assert.NoError(t, err)
	assert.Equal(t, expectedUsers, result)
	mockRepo.AssertExpectations(t)
}

func TestGetUserByID(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	user := &models.User{
		Name:     "David",
		LastName: "Marin",
		Email:    "david@example.com",
		Password: "pass123",
	}

	mockRepo.On("GetByID", mock.Anything, "abc123").Return(user, nil)

	result, err := service.GetUserByID(context.TODO(), "abc123")

	assert.NoError(t, err)
	assert.Equal(t, user, result)
	mockRepo.AssertExpectations(t)
}
