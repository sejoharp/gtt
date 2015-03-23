package controller

import (
	"net/http"

	"time"

	"encoding/json"
	"io/ioutil"

	"errors"

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
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if controller.createUser(r) != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (controller *UserControllerImpl) GetToken(c web.C, w http.ResponseWriter, r *http.Request) {
	credentials, parsingErr := parseGetTokenRequest(r)
	if parsingErr != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	user, FindUserErr := controller.userDao.FindByName(credentials.Username)
	if FindUserErr != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if controller.checkPassword(credentials) != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	expirationDate := time.Now().Add(time.Hour * 72)
	token, tokenErr := controller.tokenizer.generate(user.ID.Hex(), expirationDate)
	if tokenErr != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	response, marshalErr := json.Marshal(map[string]string{"token": token})
	if marshalErr != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Write(response)
	w.WriteHeader(http.StatusOK)
}

func parseGetTokenRequest(r *http.Request) (Credentials, error) {
	credentials := Credentials{}
	body, readErr := ioutil.ReadAll(r.Body)
	if readErr != nil {
		return credentials, readErr
	}
	parseErr := json.Unmarshal(body, &credentials)
	if parseErr != nil {
		return credentials, parseErr
	}
	return credentials, nil
}

func (controller *UserControllerImpl) checkPassword(credentials Credentials) error {
	password, getPasswordErr := controller.userDao.GetPasswordByUser(credentials.Username)
	if getPasswordErr != nil {
		return errors.New("user unknown")
	}
	hash, hashErr := controller.crypter.generateHash([]byte(credentials.Password))
	if hashErr != nil {
		return errors.New("creating hash failed")
	}
	return controller.crypter.checkPassword(hash, []byte(password))
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
