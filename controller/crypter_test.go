package controller

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Cryper", func() {
	It("should crypt a password and detects equality.", func() {
		salt := []byte("salt")
		crypter := NewCrypter(salt)
		password := []byte("password")

		passwordHash, err := crypter.generateHash(password)

		Expect(err).To(Succeed())
		completePassword := append(salt, password...)
		Expect(crypter.checkPassword(passwordHash, completePassword)).To(BeNil())
	})
})
