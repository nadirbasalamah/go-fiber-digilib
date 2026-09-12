package utils

import "golang.org/x/crypto/bcrypt"

func GeneratePassword(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword(
		[]byte(password), bcrypt.DefaultCost,
	)
}

func ComparePassword(hashed, password string) error {
	return bcrypt.CompareHashAndPassword(
		[]byte(hashed), []byte(password),
	)
}
