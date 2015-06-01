package controller

import (
	"github.com/zenazn/goji"
)

func createRouter(intervalController IntervalController, userController UserController, authController AuthController) {
	goji.Get("/token", userController.GetToken)
	goji.Post("/user", userController.Register)
	goji.Post("/interval/start", authController.authenticated(intervalController.Start))
	goji.Post("/interval/stop", authController.authenticated(intervalController.Stop))
}
