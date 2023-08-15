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
	router.Use(middleware.ClientAuth())
	
	userHandler := handler.NewUserHandler(client)
    baseRoute := router.Group("/v1/bets")

	// user endpoints
	baseRoute.POST("/user", userHandler.CreateUser)
	baseRoute.GET("/user/:userName", userHandler.GetUser)
	baseRoute.PUT("/user/:userName", userHandler.ReplaceUser)

	// bet endpoints
	baseRoute.POST("/user/bet", userHandler.CreateBet)
	baseRoute.PUT("/user/bets/:betId", userHandler.UpdateBet)
	
	router.Run()
}

