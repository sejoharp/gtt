package controller

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Cryper", func() {
	It("should crypt a password and detects equality.", func() {
		crypter := NewCrypter([]byte("salt"))
		password := []byte("password")

		passwordHash, err := crypter.generateHash(password)

		Expect(err).To(Succeed())
		completePassword := append(crypter.salt, password...)
		Expect(isSamePassword(passwordHash, completePassword)).To(BeTrue())
	})
})
