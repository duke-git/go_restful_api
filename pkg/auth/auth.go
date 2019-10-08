package auth

import (
	"golang.org/x/crypto/bcrypt"
)

//encrypts the plain text with bcrypt
func Encrypt(source string) (string, error) {
	hashedByptes, err := bcrypt.GenerateFromPassword([]byte(source), bcrypt.DefaultCost)
	return string(hashedByptes), err
}

func Compare(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}