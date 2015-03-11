package controller

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

type Tokenizer interface {
	generate(userId string) (string, error)
	parse(tokenString string) (string, error)
}

type TokenizerImpl struct {
	key []byte
}

func NewTokenizer(key []byte) Tokenizer {
	return &TokenizerImpl{key}
}

func (tokenizer *TokenizerImpl) generate(userId string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims["id"] = userId
	token.Claims["exp"] = time.Now().Add(time.Hour * 72).Unix()
	return token.SignedString(tokenizer.key)
}

func (tokenizer *TokenizerImpl) parse(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return tokenizer.key, nil
	})
	return token.Claims["id"].(string), err
}
