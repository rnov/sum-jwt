package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/qredo-external/go-rnov/pkg/user"

	"github.com/dgrijalva/jwt-go"
)

// Auth - holds all the relevant data to config JWT auth.
type Auth struct {
	secret   string
	duration time.Duration
}

// NewAuth - Auth constructor
func NewAuth(secret string, td time.Duration) *Auth {
	return &Auth{
		secret:   secret,
		duration: td,
	}
}

// Operations - defines all the business logic operations for authorization.
type Operations interface {
	CreateJWT(usr user.User) (string, error)
	ValidateJWT(JWT string) error
}

// verifyJWT - Parse, validate, and return a token.
// keyFunc will receive the parsed token and should return the key for validating.
func (a Auth) verifyJWT(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(a.secret), nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}

// ValidateJWT - validates a given JWT based on its metadata.
func (a Auth) ValidateJWT(JWT string) error {
	token, err := a.verifyJWT(JWT)
	if err != nil {
		return err
	}
	if _, ok := token.Claims.(jwt.Claims); !ok || !token.Valid {
		return errors.New("invalid token")
	}
	return nil
}

// CreateJWT - given a user create a valid JWT.
func (a Auth) CreateJWT(usr user.User) (string, error) {
	var err error
	//Creating Access Token
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["user_id"] = usr.UserName
	atClaims["exp"] = time.Now().Add(a.duration).Unix()
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, err := at.SignedString([]byte(a.secret))
	if err != nil {
		return "", err
	}
	return token, nil
}
