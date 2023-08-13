package api

import (
    "github.com/fundrick/bets-svc/internal/api/handler"
	"github.com/fundrick/bets-svc/internal/api/middleware"
	"github.com/fundrick/bets-svc/internal/bets"
    "github.com/gin-gonic/gin"

	"fmt"
)

func Start() {
	fmt.Println("Starting...")
	router := gin.Default()
	client, err := bets.CreateDynamoDBClient()
	if err != nil {
		return
	}

	router.Use(middleware.CORSMiddleware())

	userHandler := handler.NewUserHandler(client)
    baseRoute := router.Group("/v1")
	baseRoute.POST("/bets/user", userHandler.CreateUser)
	baseRoute.GET("/bets/user/:userName", userHandler.GetUser)
	baseRoute.PUT("/bets/user/:userName", userHandler.ReplaceUser)
	baseRoute.POST("bets", userHandler.CreateBet)
	baseRoute.PUT("bets/:betId", userHandler.UpdateBet)
	router.Run(":8080")
}

