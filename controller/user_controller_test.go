package controller

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"usersProject/models"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/mongo"
)

type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) CreateUser(ctx context.Context, user *models.User) (*mongo.InsertOneResult, error) {
	args := m.Called(ctx, user)
	return args.Get(0).(*mongo.InsertOneResult), args.Error(1)
}

func (m *MockUserService) GetAllUsers(ctx context.Context) ([]models.User, error) {
	args := m.Called(ctx)
	return args.Get(0).([]models.User), args.Error(1)
}

func (m *MockUserService) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserService) UpdateUser(ctx context.Context, id string, user *models.User) (*mongo.UpdateResult, error) {
	args := m.Called(ctx, id, user)
	return args.Get(0).(*mongo.UpdateResult), args.Error(1)
}

func (m *MockUserService) DeleteUser(ctx context.Context, id string) (*mongo.DeleteResult, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*mongo.DeleteResult), args.Error(1)
}

func TestCreateUserHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(MockUserService)
	userController := NewUserController(mockService)

	router := gin.Default()
	userController.RegisterRoutes(router)

	user := &models.User{
		Name:     "Juan",
		LastName: "Henao",
		Email:    "juan@example.com",
		Password: "123456",
	}

	body, _ := json.Marshal(user)
	expected := &mongo.InsertOneResult{InsertedID: "fake-id"}

	// Usamos context.TODO() como contexto de prueba
	mockService.On("CreateUser", mock.Anything, user).Return(expected, nil)

	req, _ := http.NewRequest(http.MethodPost, "/users/", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusCreated, resp.Code)
	mockService.AssertExpectations(t)
}

func TestGetAllUsersHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(MockUserService)
	controller := NewUserController(mockService)

	router := gin.Default()
	controller.RegisterRoutes(router)

	expectedUsers := []models.User{
		{Name: "Juan", LastName: "Henao", Email: "juan@example.com", Password: "123"},
		{Name: "Ana", LastName: "Gomez", Email: "ana@example.com", Password: "456"},
	}

	mockService.On("GetAllUsers", mock.Anything).Return(expectedUsers, nil)

	req, _ := http.NewRequest(http.MethodGet, "/users/", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	mockService.AssertExpectations(t)
}

func TestGetUserByIDHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(MockUserService)
	controller := NewUserController(mockService)

	router := gin.Default()
	controller.RegisterRoutes(router)

	expectedUser := &models.User{
		Name:     "Carlos",
		LastName: "López",
		Email:    "carlos@example.com",
		Password: "pass123",
	}

	mockService.On("GetUserByID", mock.Anything, "abc123").Return(expectedUser, nil)

	req, _ := http.NewRequest(http.MethodGet, "/users/abc123", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	mockService.AssertExpectations(t)
}

func TestUpdateUserHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(MockUserService)
	controller := NewUserController(mockService)

	router := gin.Default()
	controller.RegisterRoutes(router)

	updatedUser := &models.User{
		Name:     "Updated",
		LastName: "User",
		Email:    "updated@example.com",
		Password: "nuevopass",
	}

	body, _ := json.Marshal(updatedUser)

	expectedResult := &mongo.UpdateResult{MatchedCount: 1, ModifiedCount: 1}

	mockService.On("UpdateUser", mock.Anything, "abc123", updatedUser).Return(expectedResult, nil)

	req, _ := http.NewRequest(http.MethodPut, "/users/abc123", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	mockService.AssertExpectations(t)
}
func TestDeleteUserHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(MockUserService)
	controller := NewUserController(mockService)

	router := gin.Default()
	controller.RegisterRoutes(router)

	expectedResult := &mongo.DeleteResult{DeletedCount: 1}

	mockService.On("DeleteUser", mock.Anything, "abc123").Return(expectedResult, nil)

	req, _ := http.NewRequest(http.MethodDelete, "/users/abc123", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	mockService.AssertExpectations(t)
}
