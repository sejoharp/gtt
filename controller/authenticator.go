package controller

import (
	"net/http"

	"github.com/zenazn/goji/web"
)

type AuthController interface {
	authenticated(h web.HandlerFunc) web.HandlerFunc
}

type AuthControllerImpl struct {
	tokenizer Tokenizer
}

func NewAuthController(tokenizer Tokenizer) AuthController {
	return &AuthControllerImpl{tokenizer}
}
func (controller *AuthControllerImpl) authenticated(h web.HandlerFunc) web.HandlerFunc {
	return web.HandlerFunc(func(c web.C, w http.ResponseWriter, r *http.Request) {
		userID, err := controller.tokenizer.parse(r)
		if err == nil {
			c.Env["userID"] = userID
			h.ServeHTTPC(c, w, r)
		} else {
			unauthorized(c, w, r)
		}
	})
}

func unauthorized(c web.C, w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusUnauthorized)
}
