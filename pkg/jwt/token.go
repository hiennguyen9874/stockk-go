package jwt

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/hiennguyen9874/stockk-go/pkg/httpErrors"
)

type AuthClaims struct {
	jwt.StandardClaims
	Email string `json:"email"`
	Id    string `json:"id"`
}

func CreateAccessTokenHS256(id string, email string, secretKey string, expireDuration int64, issuer string) (string, error) {
	claims := AuthClaims{
		Email: email,
		Id:    id,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  time.Now().Unix(),
			Issuer:    issuer,
			ExpiresAt: time.Now().Add(time.Duration(expireDuration)).Unix(),
			Subject:   id,
			NotBefore: time.Now().Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secretKey))
}

func ParseTokenHS256(tokenString string, secretKey string) (id string, email string, err error) {
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

func CreateAccessTokenRS256(id string, email string, privateKey string, expireDuration int64, issuer string) (string, error) {
	claims := AuthClaims{
		Email: email,
		Id:    id,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  time.Now().Unix(),
			Issuer:    issuer,
			ExpiresAt: time.Now().Add(time.Duration(expireDuration)).Unix(),
			Subject:   id,
			NotBefore: time.Now().Unix(),
		},
	}

	decodedPrivateKey, err := base64.StdEncoding.DecodeString(privateKey)
	if err != nil {
		return "", err
	}
	key, err := jwt.ParseRSAPrivateKeyFromPEM(decodedPrivateKey)

	if err != nil {
		return "", err
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodRS256, claims).SignedString(key)

	if err != nil {
		return "", err
	}

	return token, nil
}

func ParseTokenRS256(tokenString string, publicKey string) (id string, email string, err error) {
	decodedPublicKey, err := base64.StdEncoding.DecodeString(publicKey)
	if err != nil {
		return "", "", err
	}

	key, err := jwt.ParseRSAPublicKeyFromPEM(decodedPublicKey)

	if err != nil {
		return "", "", err
	}

	parsedToken, err := jwt.ParseWithClaims(tokenString, &AuthClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return key, nil
	})

	if err != nil {
		return "", "", err
	}

	if !parsedToken.Valid {
		return "", "", httpErrors.Err(httpErrors.ErrorInvalidJWTToken, http.StatusBadRequest, httpErrors.ErrorInvalidJWTToken.Error())
	}

	if claims, ok := parsedToken.Claims.(*AuthClaims); ok {
		return claims.Id, claims.Email, nil
	}

	return "", "", httpErrors.Err(httpErrors.ErrorInvalidJWTClaims, http.StatusBadRequest, httpErrors.ErrorInvalidJWTClaims.Error())
}
