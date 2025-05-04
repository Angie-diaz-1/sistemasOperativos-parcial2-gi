package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type HealthController struct {
}

func NewHealthController() *HealthController {
	return &HealthController{}
}

func (uc *HealthController) RegisterRoutes(r *gin.Engine) {
	health := r.Group("/health")
	{
		health.GET("/", uc.GetHealth)
	}
}
func (uc *HealthController) GetHealth(c *gin.Context) {

	c.JSON(http.StatusOK, "Hola Docker!")
}
