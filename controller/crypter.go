package controller

import "golang.org/x/crypto/bcrypt"

type Crypter interface {
	generateHash(password []byte) ([]byte, error)
	checkPassword(hash, password []byte) error
}

type CrypterImpl struct {
	salt []byte
}

func NewCrypter(salt []byte) Crypter {
	return &CrypterImpl{salt}
}

func (crypter *CrypterImpl) generateHash(password []byte) ([]byte, error) {
	completePassword := append(crypter.salt, password...)
	return bcrypt.GenerateFromPassword(completePassword, 10)
}

func (crypter *CrypterImpl) checkPassword(hash, password []byte) error {
	return bcrypt.CompareHashAndPassword(hash, password)
}
