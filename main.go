package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"usersProject/controller"
	"usersProject/database"
	"usersProject/repository"
	"usersProject/service"
)

func main() {
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hola Docker!")
	})

	db := database.ConnectMongoDB()

	collectionName := os.Getenv("MONGO_COLLECTION_NAME")
	userRepo := repository.NewUserRepository(db, collectionName)

	userService := service.NewUserService(userRepo)
	userController := controller.NewUserController(userService)

	r := gin.Default()
	userController.RegisterRoutes(r)

	r.Run(":8080") // Corre en localhost:8080
}
