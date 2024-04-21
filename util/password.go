package util

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)



func HashPassword(password string) (string, error) {
	hashredpwd, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("Erro ao hashar, err: %w", err)
	}
	return string(hashredpwd), err
}

func CompareHashPassword(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}