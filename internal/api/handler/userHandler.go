package handler

import (
    "github.com/aws/aws-sdk-go-v2/service/dynamodb"
    "github.com/fundrick/bets-svc/internal/bets"
    "github.com/fundrick/bets-svc/internal/modals"
    "github.com/fundrick/bets-svc/internal/api/fmt"
    "net/http"
    "github.com/gin-gonic/gin"
    
)

type UserHandler struct {
    DynamoDBClient *dynamodb.Client
}

func NewUserHandler(dynamodDBClient *dynamodb.Client) *UserHandler {
    return &UserHandler{
        DynamoDBClient: dynamodDBClient,
    }
}

func (uh *UserHandler) CreateUser(c *gin.Context) {
    var newUser modals.User
    
    if err := c.BindJSON(&newUser); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to bind user data"})
        return
    }
    
    formatedUserProfile := bets.BetsUserProflieItem{
        UserName: "USER:" + newUser.UserName,
        UserAspects: "PROFILE:" + newUser.UserName,
        DiscordUserName: newUser.UserProfileAttributes.DiscordUserName,
        DiscordImageUrl: newUser.UserProfileAttributes.DiscordImageUrl,
        SteamUserName: newUser.UserProfileAttributes.SteamUserName,
        SteamUserId: newUser.UserProfileAttributes.SteamUserId,
        Wins: newUser.UserProfileAttributes.Wins,
        Losses: newUser.UserProfileAttributes.Losses,
        TopGame: newUser.UserProfileAttributes.TopGame,
        BestGame: newUser.UserProfileAttributes.BestGame,
        WorstGame: newUser.UserProfileAttributes.WorstGame,
    }
    
    if err := bets.CreateUser(uh.DynamoDBClient, formatedUserProfile); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
        return
    }

    c.IndentedJSON(http.StatusCreated, newUser)
}

func (uh *UserHandler) GetUser(c *gin.Context) {
    userName := c.Param("userName")
    formatedUsername := "USER:" + userName

    user, err := bets.GetUser(uh.DynamoDBClient, formatedUsername)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user"})
        return
    }

    formatedUser := fmt.FormatGetUserResponse(user)

    response := modals.SuccessResponse{
        Data: formatedUser,
    }

    c.IndentedJSON(http.StatusOK, response)
}

func (uh *UserHandler) GetUsers(c *gin.Context) {

    users, err := bets.GetUsers(uh.DynamoDBClient)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user"})
        return
    }

    foramtedUsers := fmt.FormatUersResponse(users)

    response := modals.SuccessResponse{
        Data: foramtedUsers,
    }

    c.IndentedJSON(http.StatusOK, response)
}