package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"usersProject/service"
)

type UserController struct {
	service service.UserServiceInterface
}

func NewUserController(s service.UserServiceInterface) *UserController {
	return &UserController{service: s}
}

func (uc *UserController) RegisterRoutes(r *gin.Engine) {
	users := r.Group("/users")
	{
		users.GET("/", uc.GetAllUsers)
		users.GET("/:id", uc.GetUserByID)

	}
}

func (uc *UserController) GetAllUsers(c *gin.Context) {
	users, err := uc.service.GetAllUsers(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, users)
}

func (uc *UserController) GetUserByID(c *gin.Context) {
	id := c.Param("id")
	user, err := uc.service.GetUserByID(c, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}
