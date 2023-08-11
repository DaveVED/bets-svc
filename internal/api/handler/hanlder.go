package handler

import (
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/dynamodb"
    "github.com/fundrick/bets-svc/internal/bets"
    "net/http"
    "github.com/gin-gonic/gin"
)

type UserHandler struct {
    DynamoSession *session.Session
    DynamoService *dynamodb.DynamoDB
}

func NewUserHandler(session *session.Session, service *dynamodb.DynamoDB) *UserHandler {
    return &UserHandler{
        DynamoSession: session,
        DynamoService: service,
    }
}

func (uh *UserHandler) CreateUser(c *gin.Context) {
    var newUser bets.NewUser
    
    if err := c.BindJSON(&newUser); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to bind user data"})
        return
    }

    /* TODO, we need to move this a func to format a newUser which we need a new TYPE for */
    formatedUser := bets.NewUser{
        UserName: newUser.UserName, // Assuming your JSON data has a field called "UserName"
        Types:  newUser.Types,  // Assuming your JSON data has a field called "Details"
    }
    
    if err := bets.CreateUser(uh.DynamoSession, uh.DynamoService, formatedUser); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
        return
    }

    c.IndentedJSON(http.StatusCreated, newUser)
}
