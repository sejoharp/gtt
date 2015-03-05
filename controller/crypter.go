package controller

import "golang.org/x/crypto/bcrypt"

type Crypter struct {
	salt []byte
}

func NewCrypter(salt []byte) *Crypter {
	return &Crypter{salt}
}

func (crypter *Crypter) generateHash(password []byte) ([]byte, error) {
	completePassword := append(crypter.salt, password...)
	return bcrypt.GenerateFromPassword(completePassword, 10)
}

func isSamePassword(hash, password []byte) bool {
	return bcrypt.CompareHashAndPassword(hash, password) == nil
}
