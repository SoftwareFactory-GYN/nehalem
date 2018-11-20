package db

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func GetSession() *dynamodb.DynamoDB {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	creds := credentials.NewEnvCredentials()

	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String("sa-east-1"),
		Credentials: creds,
	})

	// Create DynamoDB client
	return dynamodb.New(sess)

}

func ListTables(svc *dynamodb.DynamoDB) {

	result, err := svc.ListTables(&dynamodb.ListTablesInput{})

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("Tables:")
	fmt.Println("")

	for _, n := range result.TableNames {
		fmt.Println(*n)
	}

}
