package controller

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"time"

	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
	"github.com/zenazn/goji/web"
	. "github.com/zippelmann/gtt/models"
	"gopkg.in/mgo.v2/bson"
)

var _ = Describe("Configparser", func() {
	var (
		intervalController IntervalController
		responseRecorder   *httptest.ResponseRecorder
		context            web.C
		intervalDao        *IntervalDaoMock
		userID             bson.ObjectId
	)

	BeforeEach(func() {
		intervalDao = new(IntervalDaoMock)
		intervalController = NewIntervalController(intervalDao)
		responseRecorder = httptest.NewRecorder()
		context = web.C{Env: make(map[interface{}]interface{})}
		userID = bson.NewObjectId()
		context.Env["userID"] = userID
	})

	It("should start an interval.", func() {
		jsonRequest := "{}"
		httpRequest, _ := http.NewRequest("POST", "localhost", strings.NewReader(jsonRequest))
		intervalDao.On("Start", userID).Return(nil)

		intervalController.Start(context, responseRecorder, httpRequest)

		Expect(responseRecorder.Code).To(Equal(http.StatusOK))
	})

	It("should return 400, when starting an interval fails.", func() {
		jsonRequest := "{}"
		httpRequest, _ := http.NewRequest("POST", "localhost", strings.NewReader(jsonRequest))
		intervalDao.On("Start", userID).Return(errors.New("db down"))

		intervalController.Start(context, responseRecorder, httpRequest)

		Expect(responseRecorder.Code).To(Equal(http.StatusBadRequest))
	})
})

type IntervalDaoMock struct {
	mock.Mock
}

func (mock *IntervalDaoMock) Save(interval Interval) error {
	args := mock.Called(interval)
	return args.Error(0)
}
func (mock *IntervalDaoMock) FindByUserID(userID bson.ObjectId) ([]Interval, error) {
	args := mock.Called(userID)
	return args.Get(0).([]Interval), args.Error(1)
}
func (mock *IntervalDaoMock) IsUserWorking(userID bson.ObjectId) (bool, error) {
	args := mock.Called(userID)
	return args.Bool(0), args.Error(1)
}
func (mock *IntervalDaoMock) Start(userID bson.ObjectId) error {
	args := mock.Called(userID)
	return args.Error(0)
}
func (mock *IntervalDaoMock) Stop(userID bson.ObjectId) error {
	args := mock.Called(userID)
	return args.Error(0)
}
func (mock *IntervalDaoMock) FindOpenIntervals(userID bson.ObjectId) ([]Interval, error) {
	args := mock.Called(userID)
	return args.Get(0).([]Interval), args.Error(1)
}
func (mock *IntervalDaoMock) FindInRange(userID bson.ObjectId, begin time.Time, end time.Time) ([]Interval, error) {
	args := mock.Called(userID, begin, end)
	return args.Get(0).([]Interval), args.Error(1)
}
