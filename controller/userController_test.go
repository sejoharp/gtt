package controller

import (
	"net/http"
	"net/http/httptest"
	"strings"

	"errors"

	"time"

	"encoding/json"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
	"github.com/zenazn/goji/web"
	. "github.com/zippelmann/gtt/models"
	"gopkg.in/mgo.v2/bson"
)

var _ = Describe("UserController", func() {
	var (
		userDao          *UserDaoMock
		crypter          *CrypterMock
		tokenizer        *TokenizerMock
		userController   UserController
		responseRecorder *httptest.ResponseRecorder
		context          web.C
	)

	BeforeEach(func() {
		crypter = new(CrypterMock)
		userDao = new(UserDaoMock)
		tokenizer = new(TokenizerMock)
		userController = NewUserController(userDao, crypter, tokenizer)
		responseRecorder = httptest.NewRecorder()
		context = web.C{}
	})

	It("should detect valid register request.", func() {
		jsonRequest := `{"username":"peter", "password":"secret", "workTime":"2h"}`
		httpRequest, _ := http.NewRequest("GET", "localhost", strings.NewReader(jsonRequest))
		userDao.On("FindByName", mock.Anything).Return(User{}, errors.New("not existing"))
		userDao.On("SaveWithPassword", mock.Anything).Return(nil)
		crypter.On("generateHash", mock.Anything).Return([]byte("hashedPassword"), nil)

		userController.Register(context, responseRecorder, httpRequest)

		Expect(responseRecorder.Code).To(Equal(http.StatusOK))
	})

	It("should detect invalid register request - user exists.", func() {
		jsonRequest := `{"username":"peter", "password":"secret", "workTime":"2h"}`
		httpRequest, _ := http.NewRequest("GET", "localhost", strings.NewReader(jsonRequest))
		userDao.On("FindByName", mock.Anything).Return(User{}, nil)
		userDao.On("SaveWithPassword", mock.Anything).Return(nil)
		crypter.On("generateHash", mock.Anything).Return([]byte("hashedPassword"), nil)

		userController.Register(context, responseRecorder, httpRequest)

		Expect(responseRecorder.Code).To(Equal(http.StatusNotAcceptable))
	})
	It("should detect invalid register request - user too short.", func() {
		jsonRequest := `{"username":"pet", "password":"secret", "workTime":"2h"}`
		httpRequest, _ := http.NewRequest("GET", "localhost", strings.NewReader(jsonRequest))
		userDao.On("FindByName", mock.Anything).Return(User{}, errors.New("not existing"))
		userDao.On("SaveWithPassword", mock.Anything).Return(nil)
		crypter.On("generateHash", mock.Anything).Return([]byte("hashedPassword"), nil)

		userController.Register(context, responseRecorder, httpRequest)

		Expect(responseRecorder.Code).To(Equal(http.StatusNotAcceptable))
	})
	It("should detect invalid register request - invalid json.", func() {
		jsonRequest := `"username":"peter", "password":"secret", "workTime":"2h"}`
		httpRequest, _ := http.NewRequest("GET", "localhost", strings.NewReader(jsonRequest))
		userDao.On("FindByName", mock.Anything).Return(User{}, errors.New("not existing"))
		userDao.On("SaveWithPassword", mock.Anything).Return(nil)
		crypter.On("generateHash", mock.Anything).Return([]byte("hashedPassword"), nil)

		userController.Register(context, responseRecorder, httpRequest)

		Expect(responseRecorder.Code).To(Equal(http.StatusNotAcceptable))
	})
	It("should cancel registration due to a hashing error.", func() {
		jsonRequest := `{"username":"peter", "password":"secret", "workTime":"2h"}`
		httpRequest, _ := http.NewRequest("GET", "localhost", strings.NewReader(jsonRequest))
		userDao.On("FindByName", mock.Anything).Return(User{}, errors.New("not existing"))
		userDao.On("SaveWithPassword", mock.Anything).Return(nil)
		crypter.On("generateHash", mock.Anything).Return([]byte("hashedPassword"), errors.New("hashing error"))

		userController.Register(context, responseRecorder, httpRequest)

		Expect(responseRecorder.Code).To(Equal(http.StatusNotAcceptable))
	})
	It("should return an error when creating a user fails.", func() {
		jsonRequest := `{"username":"peter", "password":"secret", "workTime":"2h"}`
		httpRequest, _ := http.NewRequest("GET", "localhost", strings.NewReader(jsonRequest))
		userDao.On("FindByName", mock.Anything).Return(User{}, errors.New("not existing"))
		userDao.On("SaveWithPassword", mock.Anything).Return(errors.New("db connection failed"))
		crypter.On("generateHash", mock.Anything).Return([]byte("hashedPassword"), nil)

		userController.Register(context, responseRecorder, httpRequest)

		Expect(responseRecorder.Code).To(Equal(http.StatusNotAcceptable))
	})

	It("should return a token.", func() {
		jsonRequest := `{"username":"peter", "password":"secret"}`
		httpRequest, _ := http.NewRequest("POST", "localhost", strings.NewReader(jsonRequest))
		user := NewPersistedMinimalUser(bson.NewObjectId(), "peter", 0)
		userDao.On("FindByName", mock.Anything).Return(user, nil)
		userDao.On("GetPasswordByUser", user.Name).Return([]byte("secret"), nil)
		crypter.On("isSamePassword", mock.Anything, mock.Anything).Return(true)
		crypter.On("generateHash", mock.Anything).Return([]byte("hashedPassword"), nil)
		tokenizer.On("generate", user.ID.Hex(), mock.Anything).Return("tokenString", nil)

		userController.GetToken(context, responseRecorder, httpRequest)

		Expect(responseRecorder.Code).To(Equal(http.StatusOK))
		var responseMap map[string]string
		err := json.Unmarshal(responseRecorder.Body.Bytes(), &responseMap)
		Expect(err).To(BeNil())
		Expect(responseMap["token"]).To(Equal("tokenString"))
	})

	It("should return an error when password is wrong.", func() {
		jsonRequest := `{"username":"peter", "password":"secret"}`
		httpRequest, _ := http.NewRequest("POST", "localhost", strings.NewReader(jsonRequest))
		user := NewPersistedMinimalUser(bson.NewObjectId(), "peter", 0)
		userDao.On("FindByName", mock.Anything).Return(user, nil)
		userDao.On("GetPasswordByUser", mock.Anything).Return([]byte("realy secret"), nil)
		crypter.On("isSamePassword", mock.Anything, mock.Anything).Return(false)
		crypter.On("generateHash", mock.Anything).Return([]byte("hashedPassword"), nil)
		tokenizer.On("generate", user.ID.Hex(), mock.Anything).Return("tokenString", nil)

		userController.GetToken(context, responseRecorder, httpRequest)

		Expect(responseRecorder.Code).To(Equal(http.StatusUnauthorized))
	})

	It("should return an error when the user is unknown.", func() {
		jsonRequest := `{"username":"peter", "password":"secret"}`
		httpRequest, _ := http.NewRequest("POST", "localhost", strings.NewReader(jsonRequest))
		userDao.On("GetPasswordByUser", mock.Anything).Return([]byte(""), errors.New("no data found"))

		userController.GetToken(context, responseRecorder, httpRequest)

		Expect(responseRecorder.Code).To(Equal(http.StatusUnauthorized))
	})

	It("should return an error when generating hash fails.", func() {
		jsonRequest := `{"username":"peter", "password":"secret"}`
		httpRequest, _ := http.NewRequest("POST", "localhost", strings.NewReader(jsonRequest))
		user := NewPersistedMinimalUser(bson.NewObjectId(), "peter", 0)
		userDao.On("FindByName", mock.Anything).Return(user, nil)
		userDao.On("GetPasswordByUser", mock.Anything).Return([]byte("realy secret"), nil)
		crypter.On("isSamePassword", mock.Anything, mock.Anything).Return(true)
		crypter.On("generateHash", mock.Anything).Return([]byte("hashedPassword"), errors.New("crypting password failed"))
		tokenizer.On("generate", user.ID.Hex(), mock.Anything).Return("tokenString", nil)

		userController.GetToken(context, responseRecorder, httpRequest)

		Expect(responseRecorder.Code).To(Equal(http.StatusUnauthorized))
	})

	It("should return an error when finding the user fails.", func() {
		jsonRequest := `{"username":"peter", "password":"secret"}`
		httpRequest, _ := http.NewRequest("POST", "localhost", strings.NewReader(jsonRequest))
		user := NewPersistedMinimalUser(bson.NewObjectId(), "peter", 0)
		userDao.On("FindByName", mock.Anything).Return(user, errors.New("user unknown"))
		userDao.On("GetPasswordByUser", mock.Anything).Return([]byte("realy secret"), nil)
		crypter.On("isSamePassword", mock.Anything, mock.Anything).Return(true)
		crypter.On("generateHash", mock.Anything).Return([]byte("hashedPassword"), nil)
		tokenizer.On("generate", user.ID.Hex(), mock.Anything).Return("tokenString", nil)

		userController.GetToken(context, responseRecorder, httpRequest)

		Expect(responseRecorder.Code).To(Equal(http.StatusUnauthorized))
	})

	It("should return an error when creating token fails.", func() {
		jsonRequest := `{"username":"peter", "password":"secret"}`
		httpRequest, _ := http.NewRequest("POST", "localhost", strings.NewReader(jsonRequest))
		user := NewPersistedMinimalUser(bson.NewObjectId(), "peter", 0)
		userDao.On("FindByName", mock.Anything).Return(user, nil)
		userDao.On("GetPasswordByUser", mock.Anything).Return([]byte("realy secret"), nil)
		crypter.On("isSamePassword", mock.Anything, mock.Anything).Return(true)
		crypter.On("generateHash", mock.Anything).Return([]byte("hashedPassword"), nil)
		tokenizer.On("generate", user.ID.Hex(), mock.Anything).Return("tokenString", errors.New("creating token failed"))

		userController.GetToken(context, responseRecorder, httpRequest)

		Expect(responseRecorder.Code).To(Equal(http.StatusUnauthorized))
	})

	It("should return an error when json is invalid.", func() {
		jsonRequest := `{"username":"peter", "password":"secret"`
		httpRequest, _ := http.NewRequest("POST", "localhost", strings.NewReader(jsonRequest))
		user := NewPersistedMinimalUser(bson.NewObjectId(), "peter", 0)
		userDao.On("FindByName", mock.Anything).Return(user, nil)
		userDao.On("GetPasswordByUser", mock.Anything).Return([]byte("realy secret"), nil)
		crypter.On("isSamePassword", mock.Anything, mock.Anything).Return(true)
		crypter.On("generateHash", mock.Anything).Return([]byte("hashedPassword"), nil)
		tokenizer.On("generate", user.ID.Hex(), mock.Anything).Return(nil, nil)

		userController.GetToken(context, responseRecorder, httpRequest)

		Expect(responseRecorder.Code).To(Equal(http.StatusUnauthorized))
	})

	It("should return an error when json is invalid.", func() {
		jsonRequest := "
		";
		httpRequest, _ := http.NewRequest("POST", "localhost", strings.NewReader(jsonRequest))
		user := NewPersistedMinimalUser(bson.NewObjectId(), "peter", 0)
		userDao.On("FindByName", mock.Anything).Return(user, nil)
		userDao.On("GetPasswordByUser", mock.Anything).Return([]byte("realy secret"), nil)
		crypter.On("isSamePassword", mock.Anything, mock.Anything).Return(true)
		crypter.On("generateHash", mock.Anything).Return([]byte("hashedPassword"), nil)
		tokenizer.On("generate", user.ID.Hex(), mock.Anything).Return(nil, nil)

		userController.GetToken(context, responseRecorder, httpRequest)

		Expect(responseRecorder.Code).To(Equal(http.StatusUnauthorized))
	})
})

type CrypterMock struct {
	mock.Mock
}

func (mock CrypterMock) generateHash(password []byte) ([]byte, error) {
	args := mock.Called(password)
	return args.Get(0).([]byte), args.Error(1)
}

func (mock CrypterMock) isSamePassword(hash, password []byte) bool {
	args := mock.Called(hash, password)
	return args.Bool(0)
}

type UserDaoMock struct {
	mock.Mock
}

func (mock UserDaoMock) Save(user User) error {
	args := mock.Called(user)
	return args.Error(0)
}

func (mock UserDaoMock) SaveWithPassword(user UserWithPassword) error {
	args := mock.Called(user)
	return args.Error(0)
}

func (mock UserDaoMock) FindByID(id bson.ObjectId) (User, error) {
	args := mock.Called(id)
	return args.Get(0).(User), args.Error(1)
}

func (mock UserDaoMock) FindByName(name string) (User, error) {
	args := mock.Called(name)
	return args.Get(0).(User), args.Error(1)
}

func (mock UserDaoMock) AddPassword(id bson.ObjectId, password []byte) error {
	args := mock.Called(id, password)
	return args.Error(0)
}

func (mock UserDaoMock) AddPasswordByUser(username string, password []byte) error {
	args := mock.Called(username, password)
	return args.Error(0)
}

func (mock UserDaoMock) GetPassword(id bson.ObjectId) ([]byte, error) {
	args := mock.Called(id)
	return args.Get(0).([]byte), args.Error(1)
}

func (mock UserDaoMock) GetPasswordByUser(username string) ([]byte, error) {
	args := mock.Called(username)
	return args.Get(0).([]byte), args.Error(1)
}

func (mock UserDaoMock) Update(user User) error {
	args := mock.Called(user)
	return args.Error(0)
}

type TokenizerMock struct {
	mock.Mock
}

func (mock TokenizerMock) generate(userID string, expirationDate time.Time) (string, error) {
	args := mock.Called(userID, expirationDate)
	return args.String(0), args.Error(1)
}

func (mock TokenizerMock) parse(request *http.Request) (string, error) {
	args := mock.Called(request)
	return args.String(0), args.Error(1)
}
