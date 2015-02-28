package controller

import (
	"errors"

	"net/http"

	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("RegisterRequest", func() {
	It("should detect valid request.", func() {
		request := NewRegisterRequest("peter", "secret", "2h")
		Expect(request.validate()).To(Succeed())
	})
	It("should detect invalid password.", func() {
		request := NewRegisterRequest("peter", "sec", "2h")
		Expect(request.validate()).To(MatchError(errors.New("password too short")))
	})
	It("should detect invalid username.", func() {
		request := NewRegisterRequest("pe", "secret", "2h")
		Expect(request.validate()).To(MatchError(errors.New("username too short")))
	})
	It("should detect invalid workTime.", func() {
		request := NewRegisterRequest("peter", "secret", "")
		Expect(request.validate()).To(MatchError(errors.New("workTime not valid")))
	})
	It("should parse json request.", func() {
		jsonRequest := `{"username":"peter", "password":"secret", "workTime":"2h"}`
		httpRequest, err := http.NewRequest("GET", "localhost", strings.NewReader(jsonRequest))

		registerRequest, parseErr := parseRequest(httpRequest)

		Expect(err).To(BeNil())
		Expect(parseErr).To(BeNil())
		Expect(registerRequest.Username).To(Equal("peter"))
		Expect(registerRequest.Password).To(Equal("secret"))
		Expect(registerRequest.WorkTime).To(Equal("2h"))
	})

	It("should detect invalid json request(missing closing brace).", func() {
		jsonRequest := `{"username":"peter", "password":"secret", "workTime":"2h"`
		httpRequest, _ := http.NewRequest("GET", "localhost", strings.NewReader(jsonRequest))

		_, parseErr := parseRequest(httpRequest)

		Expect(parseErr).NotTo(BeNil())
	})

})
