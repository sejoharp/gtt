package controller

import (
	"net/http"

	"time"

	"github.com/zenazn/goji/web"
	"github.com/zippelmann/gtt/models"
)

type UserController interface {
	Register(c web.C, w http.ResponseWriter, r *http.Request)
	GetToken(c web.C, w http.ResponseWriter, r *http.Request)
}

type UserControllerImpl struct {
	userDao   models.UserDao
	crypter   Crypter
	tokenizer Tokenizer
}

func NewUserController(userDao models.UserDao, crypter Crypter, tokenizer Tokenizer) UserController {
	return &UserControllerImpl{userDao, crypter, tokenizer}
}

func (controller *UserControllerImpl) Register(c web.C, w http.ResponseWriter, r *http.Request) {
	if controller.isRegisterRequestValid(r) == false {
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}
	if saveErr := controller.createUser(r); saveErr != nil {
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (controller *UserControllerImpl) GetToken(c web.C, w http.ResponseWriter, r *http.Request) {

}

func (controller *UserControllerImpl) isRegisterRequestValid(request *http.Request) bool {
	if registerRequest, decodeErr := parseRequest(request); decodeErr != nil {
		return false
	} else if registerRequest.validate() != nil {
		return false
	} else if controller.isUsernameExisting(registerRequest.Username) {
		return false
	}
	return true
}

func (controller *UserControllerImpl) createUser(request *http.Request) error {
	parsedRequest, _ := parseRequest(request)
	workTime, _ := time.ParseDuration(parsedRequest.WorkTime)
	hash, crypterErr := controller.crypter.generateHash([]byte(parsedRequest.Password))
	if crypterErr != nil {
		return crypterErr
	}
	user := models.NewMinimalUserWithPassword(parsedRequest.Username, workTime, hash)
	saveErr := controller.userDao.SaveWithPassword(user)
	if saveErr != nil {
		return saveErr
	}
	return nil
}

func (controller *UserControllerImpl) isUsernameExisting(username string) bool {
	_, err := controller.userDao.FindByName(username)
	return err == nil
}
