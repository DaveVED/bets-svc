package api

import (
    "github.com/fundrick/bets-svc/internal/api/handler"
	"github.com/fundrick/bets-svc/internal/api/middleware"
	"github.com/fundrick/bets-svc/internal/bets"
    "github.com/gin-gonic/gin"

	"log"
)

func Start() {
	log.Println("Bets Service is Starting...")

	router := gin.Default()
	client, err := bets.CreateDynamoDBClient()
	if err != nil {
		return
	}

	router.Use(middleware.CORSMiddleware())
	//router.Use(middleware.ClientAuth())
	
	// Bets Endpoints... 
	userHandler := handler.NewUserHandler(client)
    baseRoute := router.Group("/v1/bets")

	baseRoute.GET("/health", userHandler.GetHealth)
	
	// user endpoints
	baseRoute.POST("/user", userHandler.CreateUser)
	baseRoute.GET("/users", userHandler.GetUsers)
	baseRoute.GET("/user/:userName", userHandler.GetUser)

	// bet endpoints
	baseRoute.POST("/user/bet", userHandler.CreateBet)
	baseRoute.PUT("/user/bet/:betId", userHandler.UpdateBet)
	
	router.Run()
}

