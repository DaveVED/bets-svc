package api

import (
    "github.com/fundrick/bets-svc/internal/api/handler"
	"github.com/fundrick/bets-svc/internal/bets"
    "github.com/gin-gonic/gin"

	"fmt"
)

func Start() {
	fmt.Println("Starting...")
	router := gin.Default()

	dynamoSession, dynamoService := bets.CreateSession()

	userHandler := handler.NewUserHandler(dynamoSession, dynamoService)
    baseRoute := router.Group("/v1/bets")
    baseRoute.POST("/user", userHandler.CreateUser)

	router.Run("localhost:8080")
}

