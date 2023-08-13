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
    Game string
    Amount string
	Status string
	Winner string
	CreatedOn string
	CreatedBy string
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

type UserResponse struct {
    UserName string `json:"userName"`
    UserProfileAttributes UserProfileAttributesResponse `json:"userProfileAttributesResponse"`
}

type UserProfileAttributesResponse struct {
    DiscordUserName string `json:"discordUserName"`
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
	fmt.Println(userName)

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
	
	item := items[0]
	rawUserName := items[0]["userName"].(string)
	foramtedUsername := strings.Split(rawUserName, ":")[1]
	fmt.Println(foramtedUsername)
	userResponse := UserResponse{
		UserName: foramtedUsername,
		UserProfileAttributes: UserProfileAttributesResponse {
			DiscordUserName: item["discordUsername"].(string),
			DiscordImageUrl: item["discordImageUrl"].(string),
			SteamUserName: item["steamUserName"].(string),
			SteamUserId: item["steamUserId"].(string),
			Wins: item["wins"].(string),
			Losses: item["losses"].(string),
			TopGame: item["topGame"].(string),
			BestGame: item["bestGame"].(string),
			WorstGame: item["worstGame"].(string),

		},
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
			"Game": &types.AttributeValueMemberS{Value: initiatorData.Game},
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
			"Game": &types.AttributeValueMemberS{Value: opponentData.Game},
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

func CreateDynamoDBClient() (*dynamodb.Client, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, fmt.Errorf("unable to load SKD config, %", err)
	}

	client := dynamodb.NewFromConfig(cfg)
	return client, nil
}