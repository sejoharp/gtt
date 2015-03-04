package controller

import "golang.org/x/crypto/bcrypt"

func isSamePassword(hash, password []byte) bool {
	return bcrypt.CompareHashAndPassword(hash, password) == nil
}

func generateHash(salt, password []byte) ([]byte, error) {
	completePassword := append(salt, password...)
	return bcrypt.GenerateFromPassword(completePassword, 10)
}
