package controller

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/zenazn/goji/web"
)

var _ = Describe("Authenticator", func() {
	var (
		userID               string = "54fc85502336da3e245dad5f"
		tokenizer            Tokenizer
		authController       AuthController
		expirationDatePast   time.Time = time.Date(2015, 3, 12, 8, 30, 30, 0, time.UTC)
		expirationDateFuture time.Time = time.Now().Add(time.Hour * 24)
		responseRecorder     *httptest.ResponseRecorder
		context              web.C
	)

	BeforeEach(func() {
		tokenizer = NewTokenizer([]byte("super secret key"))
		authController = NewAuthController(tokenizer)
		responseRecorder = httptest.NewRecorder()
		context = web.C{Env: make(map[interface{}]interface{})}

	})

	It("should detect a valid header and put userID into context.", func() {
		token, _ := tokenizer.generate(userID, expirationDateFuture)
		httpRequest, _ := http.NewRequest("GET", "/", nil)
		httpRequest.Header.Add("Authorization", fmt.Sprintf("Bearer %v", token))
		fakeHandler := func(c web.C, w http.ResponseWriter, r *http.Request) {
			Expect(context.Env["userID"]).To(Equal(userID))
		}

		authController.authenticated(fakeHandler).ServeHTTPC(context, responseRecorder, httpRequest)
	})
	It("should detect a invalid header.", func() {
		token, _ := tokenizer.generate(userID, expirationDatePast)
		httpRequest, _ := http.NewRequest("GET", "/", nil)
		httpRequest.Header.Add("Authorization", fmt.Sprintf("Bearer %v", token))
		fakeHandler := func(c web.C, w http.ResponseWriter, r *http.Request) {}

		authController.authenticated(fakeHandler).ServeHTTPC(context, responseRecorder, httpRequest)

		Expect(responseRecorder.Code).To(Equal(http.StatusUnauthorized))
	})

	It("should detect a missing header.", func() {
		httpRequest, _ := http.NewRequest("GET", "/", nil)
		fakeHandler := func(c web.C, w http.ResponseWriter, r *http.Request) {}

		authController.authenticated(fakeHandler).ServeHTTPC(context, responseRecorder, httpRequest)

		Expect(responseRecorder.Code).To(Equal(http.StatusUnauthorized))
	})
})
