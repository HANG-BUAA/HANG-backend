package utils

import "golang.org/x/crypto/bcrypt"

func Encrypt(stText string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(stText), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func CompareHashAndPassword(hash, pasword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(pasword))
	return err == nil
}
