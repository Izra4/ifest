package helpers

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"os"
	"time"
)

type UserClaims struct {
	ID string
	jwt.RegisteredClaims
}

func NewUserClaims(id string, exp time.Duration) UserClaims {
	return UserClaims{
		ID:               id,
		RegisteredClaims: jwt.RegisteredClaims{ExpiresAt: jwt.NewNumericDate(time.Now().Add(exp))},
	}
}

func JwtToken(id string) (string, error) {
	expStr := os.Getenv("JWT_EXP")
	exp, err := time.ParseDuration(expStr)

	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, NewUserClaims(id, exp))

	tokenStr, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "Failed to sign token", err
	}

	return tokenStr, nil
}

func DecodeJWT(signedToken string, claims *UserClaims) (string, error) {
	token, err := jwt.ParseWithClaims(signedToken, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil {
		return "", err
	}

	if !token.Valid {
		return "", errors.New("invalid token")
	}

	return signedToken, err
}
