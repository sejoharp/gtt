package controller

import (
	"time"

	"net/http"

	"fmt"

	"github.com/dgrijalva/jwt-go"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Tokenizer", func() {
	var (
		key                  []byte
		userID               string
		tokenizer            Tokenizer
		tokenString          string
		expirationDatePast   time.Time
		expirationDateFuture time.Time
	)

	BeforeEach(func() {
		key = []byte("super secret key")
		userID = "54fc85502336da3e245dad5f"
		tokenizer = NewTokenizer(key)
		tokenString = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE0MjYxNDkwMzAsImlkIjoiNTRmYzg1NTAyMzM2ZGEzZTI0NWRhZDVmIn0.l_2zai_CguQNiMaNnDDySVLLC2rFb4OPT9gnPQCSIFw"
		expirationDatePast = time.Date(2015, 3, 12, 8, 30, 30, 0, time.UTC)
		expirationDateFuture = time.Now().Add(time.Hour * 24)
	})

	It("should create a token.", func() {
		token, err := tokenizer.generate(userID, expirationDatePast)

		Expect(token).To(Equal(tokenString))
		Expect(err).To(BeNil())
	})

	It("should detect a valid token", func() {
		token, _ := tokenizer.generate(userID, expirationDateFuture)
		httpRequest, _ := http.NewRequest("GET", "/", nil)
		httpRequest.Header.Add("Authorization", fmt.Sprintf("Bearer %v", token))

		parsedUserID, err := tokenizer.parse(httpRequest)

		Expect(parsedUserID).To(Equal(userID))
		Expect(err).To(BeNil())
	})

	It("should detect expired tokens", func() {
		token, _ := tokenizer.generate(userID, expirationDatePast)
		httpRequest, _ := http.NewRequest("GET", "/", nil)
		httpRequest.Header.Add("Authorization", fmt.Sprintf("Bearer %v", token))

		parsedUserID, err := tokenizer.parse(httpRequest)

		Expect(parsedUserID).To(BeEmpty())
		Expect(err.Error()).To(Equal("token is expired"))
	})

	It("should detect a missing token", func() {
		httpRequest, _ := http.NewRequest("GET", "/", nil)

		parsedUserID, err := tokenizer.parse(httpRequest)

		Expect(parsedUserID).To(BeEmpty())
		Expect(err).To(MatchError(jwt.ErrNoTokenInRequest))
	})
})
