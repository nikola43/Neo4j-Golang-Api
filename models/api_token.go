package models

import (
	"github.com/dgrijalva/jwt-go"
)

type ApiToken struct {
	Username string
	jwt.StandardClaims
}
