package service

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/mongo"
	"testing"
	"usersProject/models"
)

type MockUserRepository struct {
	mock.Mock
}

// Implementación del método Create del repositorio
func (m *MockUserRepository) Create(ctx context.Context, user *models.User) (*mongo.InsertOneResult, error) {
	args := m.Called(ctx, user)
	result := args.Get(0)
	if result == nil {
		return nil, args.Error(1)
	}
	return result.(*mongo.InsertOneResult), args.Error(1)
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

func (m *MockUserRepository) Update(ctx context.Context, id string, user *models.User) (*mongo.UpdateResult, error) {
	args := m.Called(ctx, id, user)
	return args.Get(0).(*mongo.UpdateResult), args.Error(1)
}

func (m *MockUserRepository) Delete(ctx context.Context, id string) (*mongo.DeleteResult, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*mongo.DeleteResult), args.Error(1)
}

func TestCreateUser(t *testing.T) {
	mockRepo := new(MockUserRepository)
	userService := NewUserService(mockRepo)

	user := &models.User{
		Name:     "Test",
		LastName: "User",
		Email:    "test@example.com",
		Password: "1234",
	}

	expectedResult := &mongo.InsertOneResult{InsertedID: "fakeID123"}

	mockRepo.On("Create", mock.Anything, user).Return(expectedResult, nil)

	result, err := userService.CreateUser(context.TODO(), user)

	assert.NoError(t, err)
	assert.Equal(t, expectedResult, result)
	mockRepo.AssertExpectations(t)
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

func TestUpdateUser(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	updated := &models.User{Name: "NewName"}

	expected := &mongo.UpdateResult{MatchedCount: 1, ModifiedCount: 1}

	mockRepo.On("Update", mock.Anything, "abc123", updated).Return(expected, nil)

	result, err := service.UpdateUser(context.TODO(), "abc123", updated)

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	mockRepo.AssertExpectations(t)
}

func TestDeleteUser(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	expected := &mongo.DeleteResult{DeletedCount: 1}

	mockRepo.On("Delete", mock.Anything, "abc123").Return(expected, nil)

	result, err := service.DeleteUser(context.TODO(), "abc123")

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	mockRepo.AssertExpectations(t)
}
