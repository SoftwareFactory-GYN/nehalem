package user

import (
	"fmt"
	"github.com/SoftwareFactory-GYN/nehalem/rest_api/db"
	"github.com/SoftwareFactory-GYN/nehalem/rest_api/secret"
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

	CassandraSession := db.GetSession()

	query := fmt.Sprintf("SELECT * FROM users WHERE username='%s' ALLOW FILTERING", u.Username)
	iter := CassandraSession.Query(query).Consistency(gocql.One).Iter()
	size := iter.NumRows()

	if size > 0 {
		return true
	}

	return false
}

//  Create a new user
func (u *User) Create() error {

	CassandraSession := db.GetSession()

	// generate a unique UUID for this user
	userID := gocql.TimeUUID()

	// Salt the password
	password := HashAndSalt([]byte(u.Password))

	query := "INSERT INTO users (id, username, password) VALUES (?, ?, ?)"
	err := CassandraSession.Query(query, userID, u.Username, password).Exec()
	if err != nil {
		return err
	}

	return nil
}

// Fetch a user from the database
func FetchUser(username string) (User, error) {

	CassandraSession := db.GetSession()

	query := fmt.Sprintf("SELECT * FROM users WHERE username='%s' ALLOW FILTERING", username)
	iter := CassandraSession.Query(query).Consistency(gocql.One).Iter()

	size := iter.NumRows()

	if size == 0 {
		return User{}, fmt.Errorf("user not found")
	}

	if size > 1 {
		return User{}, fmt.Errorf("integrity error: too many users found")
	}

	m := map[string]interface{}{}
	iter.MapScan(m)

	return User{
		m["username"].(string),
		m["password"].(string),
	}, nil

}
