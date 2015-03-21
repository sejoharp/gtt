package controller

import (
	"net/http"

	"time"

	"encoding/json"
	"io/ioutil"

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
	w.WriteHeader(http.StatusOK)
}

func (controller *UserControllerImpl) GetToken(c web.C, w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	credentials := new(Credentials)
	parseErr := json.Unmarshal(body, credentials)
	if parseErr != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	password, getPasswordErr := controller.userDao.GetPasswordByUser(credentials.Username)
	if getPasswordErr != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	hash, hashErr := controller.crypter.generateHash([]byte(credentials.Password))
	if hashErr != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if controller.crypter.isSamePassword(hash, []byte(password)) == false {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	user, err := controller.userDao.FindByName(credentials.Username)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	expirationDate := time.Now().Add(time.Hour * 72)
	token, tokenErr := controller.tokenizer.generate(user.ID.Hex(), expirationDate)
	if tokenErr != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	response, marshalErr := json.Marshal(map[string]string{"token": token})
	if marshalErr != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	w.Write(response)
	w.WriteHeader(http.StatusOK)
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

type Credentials struct {
	Username string
	Password string
}
