package auth

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/vaberof/go-chat/pkg/domain"
	"strconv"
	"time"
)

var (
	ErrInvalidToken         = errors.New("token is invalid")
	ErrInvalidSigningMethod = errors.New("signing method is invalid")
	ErrExpiredToken         = errors.New("token has expired")
)

type JwtToken string

type SecretKey string

type TokenParams struct {
	UserId   domain.UserId
	TokenTtl time.Duration
}

func CreateToken(params *TokenParams, secretKey SecretKey) (string, error) {
	payload := NewPayload(params)

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    strconv.Itoa(int(payload.UserId)),
		IssuedAt:  jwt.NewNumericDate(payload.IssuedAt),
		ExpiresAt: jwt.NewNumericDate(payload.ExpiredAt),
	})

	token, err := jwtToken.SignedString([]byte(secretKey))

	return token, err
}

func VerifyToken(token string, secretKey SecretKey) (*JwtPayload, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, ErrInvalidSigningMethod
		}
		return []byte(secretKey), nil
	}

	jwtToken, err := jwt.ParseWithClaims(token, &jwt.RegisteredClaims{}, keyFunc)
	if err != nil {
		return nil, err
	}

	claims, ok := jwtToken.Claims.(*jwt.RegisteredClaims)
	if !ok {
		return nil, ErrInvalidToken
	}

	if hasExpired(claims.ExpiresAt.Time) {
		return nil, ErrExpiredToken
	}

	userId, err := strconv.Atoi(claims.Issuer)
	if err != nil {
		return nil, fmt.Errorf("failed to convert claims issuer %w", err)
	}

	payload := &JwtPayload{
		UserId:    domain.UserId(userId),
		IssuedAt:  claims.IssuedAt.Time,
		ExpiredAt: claims.ExpiresAt.Time,
	}

	return payload, nil
}

func hasExpired(expireTime time.Time) bool {
	currentTime := time.Now().UTC()
	return currentTime.After(expireTime.UTC())
}
