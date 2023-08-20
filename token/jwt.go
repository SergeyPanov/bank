package token

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

const miSecretKeySize = 32

type JWTMaker struct {
	secret string
}

func NewJWTMaker(secretKey string) (Maker, error) {
	if len(secretKey) < miSecretKeySize {
		return nil, fmt.Errorf("the secret key must be at least %d long", miSecretKeySize)
	}

	return &JWTMaker{secretKey}, nil
}

func (m *JWTMaker) CreateToken(username string, duration time.Duration) (string, *Payload, error) {
	payload, err := NewPayload(username, duration)
	if err != nil {
		return "", payload, err
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	token, err := jwtToken.SignedString([]byte(m.secret))
	return token, payload, err
}

func (m *JWTMaker) VerifyToken(token string) (*Payload, error) {
	keyF := func(t *jwt.Token) (interface{}, error) {
		_, ok := t.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, ErrInvalidToken
		}

		return []byte(m.secret), nil
	}

	jwtToken, err := jwt.ParseWithClaims(token, &Payload{}, keyF)
	if err != nil {
		verr, ok := err.(*jwt.ValidationError)
		if ok && errors.Is(verr.Inner, ErrExpiredToken) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	payload, ok := jwtToken.Claims.(*Payload)
	if !ok {
		return nil, ErrInvalidToken
	}

	return payload, nil
}
