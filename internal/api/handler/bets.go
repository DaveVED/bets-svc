package handler

import (
    "github.com/fundrick/bets-svc/internal/bets"
    "net/http"
    "github.com/google/uuid"
    "github.com/gin-gonic/gin"
    "github.com/fundrick/bets-svc/internal/modals"

    "time"
)

func (uh *UserHandler) CreateBet(c *gin.Context) {
    var newBet modals.Bet
    
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
    var newUserBet modals.BetUpdate

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

    c.IndentedJSON(http.StatusCreated, newUserBet)
}