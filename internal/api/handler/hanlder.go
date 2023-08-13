package handler

import (
    "github.com/aws/aws-sdk-go-v2/service/dynamodb"
    "github.com/fundrick/bets-svc/internal/bets"
    "net/http"
    "github.com/google/uuid"
    "github.com/gin-gonic/gin"

    "time"
)

type User struct {
    UserName string `dynamodbav:"userName"`
    UserProfileAttributes UserProfileAttributes `dynamodbav:userProfileAttributes`
}

type UserProfileAttributes struct {
    DiscordUserName string `dynamodbav:"discordUsername"`
    DiscordImageUrl string `dynamodbav:"discordImageUrl"`
    SteamUserName   string `dynamodbav:"steamUserName"`
    SteamUserId     string `dynamodav:"steamUserId"`
    Wins            string `dynamodav:"wins"`
    Losses          string `dynamodav:"losses"`
    TopGame         string `dynamodav:"topGame"`
    BestGame        string `dynamodav:"bestGame"`
    WorstGame       string `dynamodav:"worstGame"`
}

type UserBet struct {
    InitiatorUserName string `dynamodbav:"initiatorUserName"`
    OpponentUserName  string `dynamodbav:"opponentUserName"`
    Game string `dynamodbav:"game"`
    Amount string `dynamodbav:"amount"`
}

type UserHandler struct {
    DynamoDBClient *dynamodb.Client
}

func NewUserHandler(dynamodDBClient *dynamodb.Client) *UserHandler {
    return &UserHandler{
        DynamoDBClient: dynamodDBClient,
    }
}

func (uh *UserHandler) CreateUser(c *gin.Context) {
    var newUser User
    
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

    c.IndentedJSON(http.StatusOK, user)
}

func (uh *UserHandler) ReplaceUser(c *gin.Context) {
    c.IndentedJSON(http.StatusOK, "HELLO WORLD")
}

func (uh *UserHandler) CreateBet(c *gin.Context) {
    var newBet UserBet
    
    if err := c.BindJSON(&newBet); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to bind user data"})
        return
    }

    betId := uuid.NewString()
	now := time.Now()
	isoTime := now.Format(time.RFC3339)

    formatedUserBetInitiator := bets.BetItem{
        UserName: "USER:" + newBet.InitiatorUserName,
        UserAspects: "BET:" + betId,
        Opponent: newBet.OpponentUserName,
        Game: newBet.Game,
        Amount: newBet.Amount,
        Status: "IN_PROGRESS",
        Winner: "",
        CreatedOn: isoTime,
        CreatedBy: newBet.InitiatorUserName,
        UpdatedOn: isoTime,
        UpdatedBy: newBet.InitiatorUserName,
    }

    formatedUserBetOpponent := bets.BetItem{
        UserName: "USER:" + newBet.OpponentUserName,
        UserAspects: "BET:" + betId,
        Opponent: newBet.InitiatorUserName,
        Game: newBet.Game,
        Amount: newBet.Amount,
        Status: "IN_PROGRESS",
        Winner: "",
        CreatedOn: isoTime,
        CreatedBy: newBet.InitiatorUserName,
        UpdatedOn: isoTime,
        UpdatedBy: newBet.InitiatorUserName,
    }

    if err := bets.CreateBet(uh.DynamoDBClient, formatedUserBetInitiator, formatedUserBetOpponent); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
        return
    }

    c.IndentedJSON(http.StatusCreated, []bets.BetItem{formatedUserBetInitiator, formatedUserBetOpponent})
}