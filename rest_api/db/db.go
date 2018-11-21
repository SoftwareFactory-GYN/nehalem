package db

import (
	"fmt"
	"github.com/SoftwareFactory-GYN/nehalem/rest_api/utils"
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
		log.Fatal(err)
	}

	creds := credentials.NewEnvCredentials()

	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String("sa-east-1"),
		Credentials: creds,
	})

	if err != nil {
		log.Fatal(err)
	}

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

func InitTables(svc *dynamodb.DynamoDB) {

	result, err := svc.ListTables(&dynamodb.ListTablesInput{})

	if err != nil {
		log.Fatal(err)
	}

	var tables []string

	for _, n := range result.TableNames {
		tables = append(tables, string(*n))
	}

	exists, _ := utils.InArray("users", tables)

	if !exists {
		CreateUserTable(svc)
	}

}

func CreateUserTable(svc *dynamodb.DynamoDB) {

	log.Println("Creating users table...")

	input := &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("guid"),
				AttributeType: aws.String("S"),
			},
			{
				AttributeName: aws.String("username"),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("guid"),
				KeyType:       aws.String("HASH"),
			},
			{
				AttributeName: aws.String("username"),
				KeyType:       aws.String("RANGE"),
			},
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(1),
			WriteCapacityUnits: aws.Int64(1),
		},
		TableName: aws.String("users"),
	}

	_, err := svc.CreateTable(input)

	if err != nil {
		log.Fatal(err)
	}
}
