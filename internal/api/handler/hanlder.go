package handler

import (
    "github.com/aws/aws-sdk-go-v2/service/dynamodb"
    "github.com/fundrick/bets-svc/internal/bets"
    "net/http"
    "github.com/google/uuid"
    "github.com/gin-gonic/gin"

    "time"
    "fmt"
)

type User struct {
    UserName string `json:"userName"`
    UserProfileAttributes UserProfileAttributes `json:"userProfileAttributes"`
}

type UserProfileAttributes struct {
    DiscordUserName string `json:"discordUsername"`
    DiscordImageUrl string `json:"discordImageUrl"`
    SteamUserName   string `json:"steamUserName"`
    SteamUserId     string `json:"steamUserId"`
    Wins            string `json:"wins"`
    Losses          string `json:"losses"`
    TopGame         string `json:"topGame"`
    BestGame        string `json:"bestGame"`
    WorstGame       string `json:"worstGame"`
}

type UserBet struct {
    InitiatorUserName string `json:"initiatorUserName"`
    OpponentUserName  string `json:"opponentUserName"`
    Game string `json:"game"`
    Amount string `json:"amount"`
}

type UserBetUpdate struct {
    Initiator      string    `json:"initiator"`
    Opponent       string    `json:"opponent"`
    Status         string    `json:"status"`
    Accepted       string    `json:"accepted"`
    NeedsToAccept  string    `json:"needsToAccept"`
    Winner         string    `json:"winner"`
    UpdatedBy      string    `json:"updatedBy"`
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
        Initiator: newBet.InitiatorUserName,
        Game: newBet.Game,
        Amount: newBet.Amount,
        Status: "PENDING",
        Winner: "",
        NeedsToAccept: newBet.OpponentUserName,
        Accepted: "false",
        CreatedOn: isoTime,
        CreatedBy: newBet.InitiatorUserName,
        UpdatedOn: isoTime,
        UpdatedBy: newBet.InitiatorUserName,
    }

    formatedUserBetOpponent := bets.BetItem{
        UserName: "USER:" + newBet.OpponentUserName,
        UserAspects: "BET:" + betId,
        Opponent: newBet.InitiatorUserName,
        Initiator: newBet.InitiatorUserName,
        Game: newBet.Game,
        Amount: newBet.Amount,
        Status: "PENDING",
        NeedsToAccept: newBet.OpponentUserName,
        Accepted: "false",
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

func (uh *UserHandler) UpdateBet(c *gin.Context) {
    betId := c.Param("betId")
    triggeredBy := c.Request.Header.Get("triggeredBy")
    var newUserBet UserBetUpdate

    if err := c.BindJSON(&newUserBet); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to bind user data"})
        return
    }

    _, err := bets.GetUser(uh.DynamoDBClient, "USER:" + newUserBet.Initiator)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "User Initiator Does not exist."})
        return
    }

    _, err = bets.GetUser(uh.DynamoDBClient, "USER:" + newUserBet.Opponent)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "User Initiator Does not exist."})
        return
    }

    if (newUserBet.NeedsToAccept != triggeredBy) {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "You're not allowed to update this bet."})
        return
    }

    now := time.Now()
	isoTime := now.Format(time.RFC3339)
    formatedInitiatorBet := bets.BetItemUpdate{
        Winner: newUserBet.Winner,
        Accepted: newUserBet.Accepted,
        Status: newUserBet.Status,
        UpdatedOn: isoTime,
        UpdatedBy: triggeredBy,
    }
    bets.UpdateBet(uh.DynamoDBClient, "USER:" + newUserBet.Initiator, "BET:" + betId, formatedInitiatorBet)
    
    formatedOpponentBet := bets.BetItemUpdate{
        Winner: newUserBet.Winner,
        Accepted: newUserBet.Accepted,
        Status: newUserBet.Status,
        UpdatedOn: isoTime,
        UpdatedBy: triggeredBy,
    }
    bets.UpdateBet(uh.DynamoDBClient, "USER:" + newUserBet.Opponent, "BET:" + betId, formatedOpponentBet)

    fmt.Println(formatedOpponentBet)
    c.IndentedJSON(http.StatusCreated, newUserBet)
}