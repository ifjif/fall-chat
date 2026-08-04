package utils

import "golang.org/x/crypto/bcrypt"

func EncryptPassword(password string) string {
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), 12)

	return string(hash)
}

func DecryptPassword(hash string, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err
}
