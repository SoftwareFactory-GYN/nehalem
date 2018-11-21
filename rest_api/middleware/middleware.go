package middleware

import (
	"github.com/SoftwareFactory-GYN/nehalem/rest_api/secret"
	"github.com/auth0/go-jwt-middleware"
	"github.com/dgrijalva/jwt-go"
)

var JwtMiddleware = jwtmiddleware.New(jwtmiddleware.Options{
	ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
		return secret.GetSigningKey(), nil
	},
	SigningMethod: jwt.SigningMethodHS256,
})
