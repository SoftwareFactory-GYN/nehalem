package user

import (
	"fmt"
	"github.com/SoftwareFactory-GYN/nehalem/rest_api/db"
	"github.com/SoftwareFactory-GYN/nehalem/rest_api/secret"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"github.com/dgrijalva/jwt-go"
	"github.com/gocql/gocql"
	"golang.org/x/crypto/bcrypt"
	"log"
	"time"
)

type User struct {
	Guid     string `json:"guid,omitempty"`
	Username string `json:"username"`
	Password string `json:"password"`
}
type UserToken struct {
	Token string `json:"token"`
}

// Salt and hash
func HashAndSalt(pwd []byte) string {

	// Use GenerateFromPassword to hash & salt pwd.
	// MinCost is just an integer constant provided by the bcrypt
	// package along with DefaultCost & MaxCost.
	// The cost can be any value you want provided it isn't lower
	// than the MinCost (4)
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}

	// GenerateFromPassword returns a byte slice so we need to
	// convert the bytes to a string and return it
	return string(hash)
}

// Compare salts
func ComparePasswords(hashedPwd string, plainPwd []byte) bool {

	// Since we'll be getting the hashed password from the DB it
	// will be a string so we'll need to convert it to a byte slice
	byteHash := []byte(hashedPwd)

	err := bcrypt.CompareHashAndPassword(byteHash, plainPwd)
	if err != nil {
		log.Println(err)
		return false
	}

	return true

}

// Create a token for a user
func (u *User) GetToken() string {
	/* Create the token */
	token := jwt.New(jwt.SigningMethodHS256)

	//Create a map to store our claims
	claims := token.Claims.(jwt.MapClaims)

	/* Set token claims */
	claims["name"] = u.Username
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	/* Sign the token with our secret */
	tokenString, _ := token.SignedString(secret.GetSigningKey())

	return tokenString
}

// Check if user exists in persistence
func (u *User) Exists() bool {

	svc := db.GetSession()

	filt := expression.Name("username").Equal(expression.Value(u.Username))

	// Get back the title, year, and rating
	proj := expression.NamesList(expression.Name("username"), expression.Name("guid"), expression.Name("password"))

	expr, err := expression.NewBuilder().WithFilter(filt).WithProjection(proj).Build()

	if err != nil {
		log.Fatal(err)
	}

	params := &dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		ProjectionExpression:      expr.Projection(),
		TableName:                 aws.String("users"),
	}

	// Make the DynamoDB Query API call
	result, err := svc.Scan(params)

	if err != nil {
		log.Fatal(err)
	}

	exists := *result.Count

	if exists > 0 {
		return true
	}

	return false
}

//  Create a new user
func (u *User) Create() error {

	svc := db.GetSession()

	// generate a unique UUID for this user
	userID, _ := gocql.RandomUUID()

	u.Guid = userID.String()

	// Salt the password
	u.Password = HashAndSalt([]byte(u.Password))

	av, err := dynamodbattribute.MarshalMap(u)
	if err != nil {
		return err
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String("users"),
	}

	_, err = svc.PutItem(input)

	return err
}

// Fetch a user from the database
func FetchUser(username string) (User, error) {

	svc := db.GetSession()

	filt := expression.Name("username").Equal(expression.Value(username))

	// Get back the title, year, and rating
	proj := expression.NamesList(expression.Name("username"), expression.Name("guid"), expression.Name("password"))

	expr, err := expression.NewBuilder().WithFilter(filt).WithProjection(proj).Build()

	if err != nil {
		log.Fatal(err)
	}

	params := &dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		ProjectionExpression:      expr.Projection(),
		TableName:                 aws.String("users"),
	}

	// Make the DynamoDB Query API call
	result, err := svc.Scan(params)

	if err != nil {
		log.Fatal(err)
	}

	exists := *result.Count

	if exists == 0 {
		return User{}, fmt.Errorf("integrity error: user not found")
	} else if exists > 1 {
		return User{}, fmt.Errorf("integrity error: more than one user found")
	}

	user := User{}

	for _, i := range result.Items {

		err = dynamodbattribute.UnmarshalMap(i, &user)

		if err != nil {
			log.Fatal(err)
		}

		break

	}

	return user, nil

}
