package main

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gocql/gocql"
	"golang.org/x/crypto/bcrypt"
	"log"
	"time"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
type UserToken struct {
	Token string `json:"token"`
}

// Salt and hash
func hashAndSalt(pwd []byte) string {

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

func (u *User) getToken() string {
	/* Create the token */
	token := jwt.New(jwt.SigningMethodHS256)

	//Create a map to store our claims
	claims := token.Claims.(jwt.MapClaims)

	/* Set token claims */
	claims["name"] = u.Username
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	/* Sign the token with our secret */
	tokenString, _ := token.SignedString(mySigningKey)

	return tokenString
}

func (u *User) exists() bool {

	query := fmt.Sprintf("SELECT * FROM users WHERE username='%s' ALLOW FILTERING", u.Username)
	iter := CassandraSession.Query(query).Consistency(gocql.One).Iter()
	size := iter.NumRows()

	if size > 1 {
		return true
	}

	return false
}

func (u *User) create() error {

	// generate a unique UUID for this user
	userID := gocql.TimeUUID()

	// Salt the password
	password := hashAndSalt([]byte(u.Password))

	query := "INSERT INTO users (id, username, password) VALUES (?, ?, ?)"
	err := CassandraSession.Query(query, userID, u.Username, password).Exec()
	if err != nil {
		return err
	}

	return nil
}
