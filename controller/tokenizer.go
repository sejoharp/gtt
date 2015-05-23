package controller

import (
	"time"

	"net/http"

	"github.com/dgrijalva/jwt-go"
)

type Tokenizer interface {
	generate(userId string, expirationDate time.Time) (string, error)
	parse(request *http.Request) (string, error)
}

type TokenizerImpl struct {
	key []byte
}

func NewTokenizer(key []byte) Tokenizer {
	return &TokenizerImpl{key}
}

func (tokenizer *TokenizerImpl) generate(userId string, expirationDate time.Time) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims["id"] = userId
	token.Claims["exp"] = expirationDate.Unix()
	return token.SignedString(tokenizer.key)
}

func (tokenizer *TokenizerImpl) parse(request *http.Request) (string, error) {
	token, err := jwt.ParseFromRequest(request, tokenizer.getTokenKey)
	if err == nil {
		return token.Claims["id"].(string), err
	}
	return "", err

}

func (tokenizer *TokenizerImpl) getTokenKey(token *jwt.Token) (interface{}, error) {
	return tokenizer.key, nil
}
