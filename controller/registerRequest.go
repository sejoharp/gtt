package controller

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

type RegisterRequest struct {
	Username string
	Password string
	WorkTime string
}

func NewRegisterRequest(username, password, workTime string) RegisterRequest {
	return RegisterRequest{username, password, workTime}
}

func parseRequest(request *http.Request) (RegisterRequest, error) {
	var parsedRequest RegisterRequest
	decoder := json.NewDecoder(request.Body)
	deCodeErr := decoder.Decode(&parsedRequest)
	return parsedRequest, deCodeErr
}

func (request RegisterRequest) validate() error {
	if len(request.Username) < 4 {
		return errors.New("username too short")
	}
	if len(request.Password) < 6 {
		return errors.New("password too short")
	}
	_, parseErr := time.ParseDuration(request.WorkTime)
	if parseErr != nil {
		return errors.New("workTime not valid")
	}
	return nil
}
