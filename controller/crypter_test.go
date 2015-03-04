package controller

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Cryper", func() {
	It("should crypt a password and detects equality.", func() {
		passwordHash, err := generateHash([]byte("salt"), []byte("password"))

		Expect(err).To(Succeed())
		completePassword := append([]byte("salt"), []byte("password")...)
		Expect(isSamePassword(passwordHash, completePassword)).To(BeTrue())
	})
})
