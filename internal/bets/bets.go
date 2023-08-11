package bets

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type UserItem struct {
	UserName string `dynamodbav:"user"`
	Types  string `dynamodbav:"types"`
    Details string `dynamodbav:"details"`
}

type DiscordType struct {
	DiscordUsername  string `json:"discordUsername"`
	DiscordImageUrl string `json:"discordImageUrl"`
}

type UserType struct {
	DISCORD DiscordType `json:"DISCORD"`
}

type NewUser struct {
	UserName string    `json:"userName"`
	Types    []UserType `json:"types"`
}

func CreateUser(sess *session.Session, svc *dynamodb.DynamoDB, data NewUser) error {
	item := UserItem{
		UserName: "USER:" + data.UserName,
		Types:  "DISCORD:" + data.UserName,
        Details: `{
            "discordUsername": "daveved12",
            "discordImageUrl": "www.abc.com"
        }`,
	}

    fmt.Println(item)
	av, err := dynamodbattribute.MarshalMap(item)
	if err != nil {
		return fmt.Errorf("failed to marshal new user item: %w", err)
	}

	tableName := "bets-dev-table"

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(tableName),
	}

	_, err = svc.PutItem(input)
	if err != nil {
        fmt.Println("HERE")
        fmt.Println(err)
		return fmt.Errorf("failed to call PutItem: %w", err)
	}

	return nil
}

/*func GetUser(sess *session.Session, svc *dynamodb.DynamoDB) {
	result, err := svc.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String("bets-dev-table"),
		Key: map[string]*dynamodb.AttributeValue{
			"user": {
				S: aws.String("#USER:TEST1"),
			},
			"details": {
				S: aws.String("PROFILE:TEST1"),
			},
		},
	})

	if err != nil {
		fmt.Println("GetItem API call failed:")
		fmt.Println(err.Error())
		return
	}

	item := UserItem{}

	err = dynamodbattribute.UnmarshalMap(result.Item, &item)

	if err != nil {
		panic(fmt.Sprintf("Failed to unmarshal Record, %v", err))
	}

	if item.UserName == "" {
		fmt.Println("Could not find Item")
		return
	}

	fmt.Println("Found item:")
	fmt.Println("User:  ", item.UserName)
	fmt.Println("Details: ", item.Details)
}*/

func CreateSession() (*session.Session, *dynamodb.DynamoDB) {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	svc := dynamodb.New(sess)

	return sess, svc
}
