package jwt

import (
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/hiennguyen9874/stockk-go/pkg/httpErrors"
)

type AuthClaims struct {
	jwt.StandardClaims
	Email string `json:"email"`
	Id    string `json:"id"`
}

func CreateAccessToken(id string, email string, secretKey string, expireDuration int64, issuer string) (string, error) {
	claims := AuthClaims{
		Email: email,
		Id:    id,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  time.Now().Unix(),
			Issuer:    issuer,
			ExpiresAt: time.Now().Add(time.Duration(expireDuration)).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secretKey))
}

func ParseToken(tokenString string, secretKey string) (id string, email string, err error) {
	token, err := jwt.ParseWithClaims(tokenString, &AuthClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})

	if err != nil {
		return "", "", err
	}

	if !token.Valid {
		return "", "", httpErrors.Err(httpErrors.ErrorInvalidJWTToken, http.StatusBadRequest, httpErrors.ErrorInvalidJWTToken.Error())
	}

	if claims, ok := token.Claims.(*AuthClaims); ok {
		return claims.Id, claims.Email, nil
	}

	return "", "", httpErrors.Err(httpErrors.ErrorInvalidJWTClaims, http.StatusBadRequest, httpErrors.ErrorInvalidJWTClaims.Error())
}
