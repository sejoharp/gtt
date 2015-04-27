package controller

import (
	"net/http"

	"github.com/zenazn/goji/web"
	"github.com/zippelmann/gtt/models"
)

type IntervalController interface {
	Start(c web.C, w http.ResponseWriter, r *http.Request)
	Stop(c web.C, w http.ResponseWriter, r *http.Request)
}

func NewIntervalController(intervalDao models.IntervalDao) IntervalController {
	return &IntervalControllerImpl{intervalDao}
}

type IntervalControllerImpl struct {
	userDao models.IntervalDao
}

func (controller *IntervalControllerImpl) Start(c web.C, w http.ResponseWriter, r *http.Request) {
	err := controller.userDao.Start(getUserID(c))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
	w.WriteHeader(http.StatusOK)
}

func (controller *IntervalControllerImpl) Stop(c web.C, w http.ResponseWriter, r *http.Request) {
	err := controller.userDao.Stop(getUserID(c))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
	w.WriteHeader(http.StatusOK)
}
