package bets

import (
	"context"
	"fmt"
	"strings"
	
    "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
    "github.com/aws/aws-sdk-go-v2/service/dynamodb"
    "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
)

type BetItem struct {
	UserName string 
    UserAspects string
	Opponent string
    Initiator string
    Game string
    Accepted string
    Amount string
    NeedsToAccept string
	Status string
	Winner string
	CreatedOn string
	CreatedBy string
	UpdatedOn string
	UpdatedBy string
}

type BetItemUpdate struct {
    Game string
    Accepted string
	Status string
	Winner string
	UpdatedOn string
	UpdatedBy string
}

type BetsUserProflieItem struct {
	UserName string
	UserAspects  string
    DiscordUserName string
    DiscordImageUrl string
    SteamUserName   string
    SteamUserId     string
	Wins            string
    Losses          string
    TopGame         string
    BestGame        string
    WorstGame       string
}


type BetResponse struct {
    BetId     string `json:"betId"`
    Initiator string `json:"initiator"`
    Opponent  string `json:"opponent"`
    Game      string `json:"game"`
    Amount    string `json:"amount"`
    Status    string `json:"status"`
    Accepted  string `json:"accepted"`
    NeedsToAccept string `json:"needsToAccept"`
    Winner    string `json:"winner"`
    CreatedOn string `json:"createdOn"`
    CreatedBy string `json:"createdBy"`
    UpdatedOn string `json:"updatedOn"`
    UpdatedBy string `json:"updatedBy"`
}

type UserResponse struct {
    UserName string `json:"userName"`
    UserProfileAttributes UserProfileAttributesResponse `json:"userProfileAttributesResponse,omitempty"`
    Bets []BetResponse `json:"bets,omitempty"`
}

type UserProfileAttributesResponse struct {
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

func CreateUser(svc *dynamodb.Client, data BetsUserProflieItem) error {
    _, err := svc.PutItem(context.TODO(), &dynamodb.PutItemInput{
        TableName: aws.String("bets-dev-table"),
        Item: map[string]types.AttributeValue{
            "userName":    &types.AttributeValueMemberS{Value: data.UserName},
            "userAspects":  &types.AttributeValueMemberS{Value: data.UserAspects},
            "discordUsername": &types.AttributeValueMemberS{Value: data.DiscordUserName},
			"discordImageUrl": &types.AttributeValueMemberS{Value: data.DiscordImageUrl},
            "steamUserName": &types.AttributeValueMemberS{Value: data.SteamUserName},
            "steamUserId": &types.AttributeValueMemberS{Value: data.SteamUserId},
            "wins": &types.AttributeValueMemberS{Value: data.Wins},
			"losses": &types.AttributeValueMemberS{Value: data.Losses},
            "topGame": &types.AttributeValueMemberS{Value: data.TopGame},
            "bestGame": &types.AttributeValueMemberS{Value: data.BestGame},
            "worstGame": &types.AttributeValueMemberS{Value: data.WorstGame},
        },
    })

	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

func GetUser(svc *dynamodb.Client, userName string) (UserResponse, error) {
	var userResponse UserResponse

	res, err := svc.Query(context.TODO(), &dynamodb.QueryInput{
		TableName: aws.String("bets-dev-table"),
        KeyConditionExpression: aws.String("userName = :hashKey"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":hashKey":  &types.AttributeValueMemberS{Value: userName},
		},
	})

	if err != nil {
		return UserResponse{}, fmt.Errorf("failed to get user: %w", err)
	}
	
	var items []map[string]interface{}
	err = attributevalue.UnmarshalListOfMaps(res.Items, &items)
	if err != nil {
		return UserResponse{}, fmt.Errorf("failed to unmarshall user: %w", err)
	}
	
    for _, v := range items {
        userAspect := v["userAspects"].(string)
        splitAspectType := strings.Split(userAspect, ":")
        aspectType := splitAspectType[0]
        aspectValue := splitAspectType[1]
        switch aspectType {
        case "PROFILE":
            userResponse.UserName = strings.Split(v["userName"].(string), ":")[1]
            userResponse.UserProfileAttributes = UserProfileAttributesResponse{
                DiscordUserName: v["discordUsername"].(string),
                DiscordImageUrl: v["discordImageUrl"].(string),
                SteamUserName: v["steamUserName"].(string),
                SteamUserId: v["steamUserId"].(string),
                Wins: v["wins"].(string),
                Losses: v["losses"].(string),
                TopGame: v["topGame"].(string),
                BestGame: v["bestGame"].(string),
                WorstGame: v["worstGame"].(string),
            }
        case "BET":
            fmt.Println(aspectValue)
            bet := BetResponse{
                BetId:     aspectValue,
                Opponent:  v["opponent"].(string),
                Game:      v["game"].(string),
                Amount:    v["amount"].(string),
                Status:    v["status"].(string),
                Winner:    v["winner"].(string),
                NeedsToAccept: v["needsToAccept"].(string),
                Initiator: v["initiator"].(string),
                Accepted: v["accepted"].(string),
                CreatedOn: v["createdOn"].(string),
                CreatedBy: v["createdBy"].(string),
                UpdatedOn: v["updatedOn"].(string),
                UpdatedBy: v["updatedBy"].(string),
            }
            userResponse.Bets = append(userResponse.Bets, bet)
        }
    }

	return userResponse, nil
}

func UpdateUser(svc *dynamodb.Client, data BetsUserProflieItem) error {
	fmt.Println(data.UserName)
	fmt.Println(data.UserAspects)

	key := map[string]types.AttributeValue{
		"userName":    &types.AttributeValueMemberS{Value: data.UserName},
		"userAspects": &types.AttributeValueMemberS{Value: data.UserAspects},
	}	

    // Modify the update expression to exclude userAspects
    updateExpression := "SET discordUsername = :discordUsername, discordImageUrl = :discordImageUrl, steamUserName = :steamUserName, steamUserId = :steamUserId, wins = :wins, losses = :losses, topGame = :topGame, bestGame = :bestGame, worstGame = :worstGame"
    
    attributeValues := map[string]types.AttributeValue{
        ":discordUsername": &types.AttributeValueMemberS{Value: data.DiscordUserName},
        ":discordImageUrl": &types.AttributeValueMemberS{Value: data.DiscordImageUrl},
        ":steamUserName": &types.AttributeValueMemberS{Value: data.SteamUserName},
        ":steamUserId": &types.AttributeValueMemberS{Value: data.SteamUserId},
        ":wins": &types.AttributeValueMemberS{Value: data.Wins},
        ":losses": &types.AttributeValueMemberS{Value: data.Losses},
        ":topGame": &types.AttributeValueMemberS{Value: data.TopGame},
        ":bestGame": &types.AttributeValueMemberS{Value: data.BestGame},
        ":worstGame": &types.AttributeValueMemberS{Value: data.WorstGame},
    }

    _, err := svc.UpdateItem(context.TODO(), &dynamodb.UpdateItemInput{
        TableName:        aws.String("bets-dev-table"),
        Key:              key,
        UpdateExpression: aws.String(updateExpression),
        ExpressionAttributeValues: attributeValues,
    })

    if err != nil {
		fmt.Println(err)
        return fmt.Errorf("failed to update user: %w", err)
    }

    return nil
}

func CreateBet(svc *dynamodb.Client, initiatorData BetItem, opponentData BetItem) error {
    _, err1 := svc.PutItem(context.TODO(), &dynamodb.PutItemInput{
        TableName: aws.String("bets-dev-table"),
        Item: map[string]types.AttributeValue{
            "userName":    &types.AttributeValueMemberS{Value: initiatorData.UserName},
            "userAspects":  &types.AttributeValueMemberS{Value: initiatorData.UserAspects},
            "opponent": &types.AttributeValueMemberS{Value: initiatorData.Opponent},
            "initiator": &types.AttributeValueMemberS{Value: initiatorData.Initiator},
            "accepted": &types.AttributeValueMemberS{Value: initiatorData.Accepted},
            "needsToAccept": &types.AttributeValueMemberS{Value: initiatorData.NeedsToAccept},
			"game": &types.AttributeValueMemberS{Value: initiatorData.Game},
			"amount": &types.AttributeValueMemberS{Value: initiatorData.Amount},
			"winner": &types.AttributeValueMemberS{Value: opponentData.Winner},
			"status": &types.AttributeValueMemberS{Value: initiatorData.Status},
            "createdOn": &types.AttributeValueMemberS{Value: initiatorData.CreatedOn},
            "createdBy": &types.AttributeValueMemberS{Value: initiatorData.CreatedBy},
			"updatedOn": &types.AttributeValueMemberS{Value: initiatorData.UpdatedOn},
            "updatedBy": &types.AttributeValueMemberS{Value: initiatorData.UpdatedBy},
        },
    })

	if err1 != nil {
		return fmt.Errorf("failed to create user: %w", err1)
	}

	_, err2 := svc.PutItem(context.TODO(), &dynamodb.PutItemInput{
        TableName: aws.String("bets-dev-table"),
        Item: map[string]types.AttributeValue{
            "userName":    &types.AttributeValueMemberS{Value: opponentData.UserName},
            "userAspects":  &types.AttributeValueMemberS{Value: opponentData.UserAspects},
            "opponent": &types.AttributeValueMemberS{Value: opponentData.Opponent},
			"game": &types.AttributeValueMemberS{Value: opponentData.Game},
			"amount": &types.AttributeValueMemberS{Value: initiatorData.Amount},
            "initiator": &types.AttributeValueMemberS{Value: initiatorData.Initiator},
            "accepted": &types.AttributeValueMemberS{Value: initiatorData.Accepted},
            "needsToAccept": &types.AttributeValueMemberS{Value: initiatorData.NeedsToAccept},
			"winner": &types.AttributeValueMemberS{Value: opponentData.Winner},
			"status": &types.AttributeValueMemberS{Value: opponentData.Status},
            "createdOn": &types.AttributeValueMemberS{Value: opponentData.CreatedOn},
            "createdBy": &types.AttributeValueMemberS{Value: opponentData.CreatedBy},
			"updatedOn": &types.AttributeValueMemberS{Value: opponentData.UpdatedOn},
            "updatedBy": &types.AttributeValueMemberS{Value: opponentData.UpdatedBy},
        },
    })

	if err2 != nil {
		return fmt.Errorf("failed to create user: %w", err2)
	}
	return nil
}


func UpdateBet(svc *dynamodb.Client, userName string, betId string, data BetItemUpdate) error {
	key := map[string]types.AttributeValue{
		"userName":    &types.AttributeValueMemberS{Value: userName},
		"userAspects": &types.AttributeValueMemberS{Value: betId},
	}	

    // Modify the update expression to exclude userAspects
    updateExpression := "SET winner = :winner, #Status = :status, accepted = :accepted, updatedBy = :updatedBy, updatedOn = :updatedOn"
    
    attributeValues := map[string]types.AttributeValue{
        ":winner": &types.AttributeValueMemberS{Value: data.Winner},
        ":status": &types.AttributeValueMemberS{Value: data.Status},
        ":accepted": &types.AttributeValueMemberS{Value: data.Accepted},
        ":updatedBy": &types.AttributeValueMemberS{Value: data.UpdatedBy},
        ":updatedOn": &types.AttributeValueMemberS{Value: data.UpdatedOn},
    }

    _, err := svc.UpdateItem(context.TODO(), &dynamodb.UpdateItemInput{
        TableName:        aws.String("bets-dev-table"),
        Key:              key,
        UpdateExpression: aws.String(updateExpression),
        ExpressionAttributeNames: map[string]string{
            "#Status": "status",
        },
        ExpressionAttributeValues: attributeValues,
    })

    if err != nil {
        return fmt.Errorf("failed to update user: %w", err)
    }

    return nil
}


func CreateDynamoDBClient() (*dynamodb.Client, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, fmt.Errorf("unable to load SKD config, %", err)
	}

	client := dynamodb.NewFromConfig(cfg)
	return client, nil
}