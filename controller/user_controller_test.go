package controller

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
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

func TestCreateUserHandler_BadRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(MockUserService)
	controller := NewUserController(mockService)
	router := gin.Default()
	controller.RegisterRoutes(router)

	// Se envía un JSON inválido
	invalidJSON := []byte(`{"name": "Juan", "email": "juan@example.com", "password": "123456"`) // Falta cierre de llave

	req, _ := http.NewRequest(http.MethodPost, "/users/", bytes.NewBuffer(invalidJSON))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	// Se espera HTTP 400 por error de binding
	assert.Equal(t, http.StatusBadRequest, resp.Code)
	// No se invoca el servicio si falla el binding, por lo que no se registran expectativas.
}

func TestCreateUserHandler_ServiceError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(MockUserService)
	controller := NewUserController(mockService)
	router := gin.Default()
	controller.RegisterRoutes(router)

	user := &models.User{
		Name:     "Juan",
		LastName: "Henao",
		Email:    "juan@example.com",
		Password: "123456",
	}
	body, _ := json.Marshal(user)

	// Se simula un error en el servicio
	mockService.On("CreateUser", mock.Anything, user).Return((*mongo.InsertOneResult)(nil), errors.New("error en la creación"))

	req, _ := http.NewRequest(http.MethodPost, "/users/", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusInternalServerError, resp.Code)
	mockService.AssertExpectations(t)
}

func TestGetAllUsersHandler_ServiceError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(MockUserService)
	controller := NewUserController(mockService)
	router := gin.Default()
	controller.RegisterRoutes(router)

	// Simula error al obtener los usuarios
	mockService.On("GetAllUsers", mock.Anything).Return(nil, errors.New("error al obtener usuarios"))

	req, _ := http.NewRequest(http.MethodGet, "/users/", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusInternalServerError, resp.Code)
	mockService.AssertExpectations(t)
}

func TestGetUserByIDHandler_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(MockUserService)
	controller := NewUserController(mockService)
	router := gin.Default()
	controller.RegisterRoutes(router)

	// Se simula que el usuario no existe devolviendo un error
	mockService.On("GetUserByID", mock.Anything, "nonexistent").Return((*models.User)(nil), errors.New("usuario no encontrado"))

	req, _ := http.NewRequest(http.MethodGet, "/users/nonexistent", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	// Según el controlador se retorna 404 si hay error en GetUserByID
	assert.Equal(t, http.StatusNotFound, resp.Code)
	mockService.AssertExpectations(t)
}

func TestUpdateUserHandler_BadRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(MockUserService)
	controller := NewUserController(mockService)
	router := gin.Default()
	controller.RegisterRoutes(router)

	// Se envía JSON mal formado
	invalidJSON := []byte(`{"name": "Updated", "email": "updated@example.com"`) // Falta campos o cierre

	req, _ := http.NewRequest(http.MethodPut, "/users/abc123", bytes.NewBuffer(invalidJSON))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
}

func TestUpdateUserHandler_ServiceError(t *testing.T) {
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

	// Se simula un error en la actualización
	mockService.On("UpdateUser", mock.Anything, "abc123", updatedUser).Return((*mongo.UpdateResult)(nil), errors.New("error al actualizar"))

	req, _ := http.NewRequest(http.MethodPut, "/users/abc123", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusInternalServerError, resp.Code)
	mockService.AssertExpectations(t)
}

func TestDeleteUserHandler_ServiceError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(MockUserService)
	controller := NewUserController(mockService)
	router := gin.Default()
	controller.RegisterRoutes(router)

	// Se simula un error al eliminar el usuario
	mockService.On("DeleteUser", mock.Anything, "abc123").Return((*mongo.DeleteResult)(nil), errors.New("error al eliminar"))

	req, _ := http.NewRequest(http.MethodDelete, "/users/abc123", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusInternalServerError, resp.Code)
	mockService.AssertExpectations(t)
}
