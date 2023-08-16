package middleware

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/gin-gonic/gin"
)

type SecretData struct {
	ClientId     string `json:"CLIENT_ID"`
	ClientSecret string `json:"CLIENT_SECRET"`
}

func ClientAuth() gin.HandlerFunc {
	return func(c *gin.Context) {

		if c.FullPath() == "/v1/bets/health" {
			c.Next()
			return
		}

		secretName := "clients/creds/bets-ui"
		region := "us-east-1"

		cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
		if err != nil {
			log.Println("Error loading AWS config:", err)
			c.String(http.StatusInternalServerError, "Internal Server Error")
			c.Abort()
			return
		}

		svc := secretsmanager.NewFromConfig(cfg)

		input := &secretsmanager.GetSecretValueInput{
			SecretId:     aws.String(secretName),
			VersionStage: aws.String("AWSCURRENT"),
		}

		result, err := svc.GetSecretValue(context.TODO(), input)
		if err != nil {
			log.Println("Error fetching secret:", err)
			c.String(http.StatusInternalServerError, "Internal Server Error")
			c.Abort()
			return
		}

		var secretString string = *result.SecretString

		var secretData SecretData
		err = json.Unmarshal([]byte(secretString), &secretData)
		if err != nil {
			log.Println("Error unmarshalling secret:", err)
			c.String(http.StatusInternalServerError, "Internal Server Error")
			c.Abort()
			return
		}

		clientId := c.Request.Header.Get("clientId")
		clientSecret := c.Request.Header.Get("clientSecret") // Fixed typo here

		if clientId != secretData.ClientId || clientSecret != secretData.ClientSecret {
			c.String(http.StatusUnauthorized, "Unauthorized")
			c.Abort()
			return
		}

		c.Next()
	}
}