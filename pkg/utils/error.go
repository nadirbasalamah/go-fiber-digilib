package utils

import (
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func GetErrorMessage(err error, data ...string) string {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return "record not found"
	} else if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		return "password not match"
	} else if errors.Is(err, gorm.ErrDuplicatedKey) {
		if len(data) > 0 {
			return fmt.Sprintf("%s must be unique", data[0])
		}
		return "data must be unique"
	} else if errors.Is(err, gorm.ErrInvalidTransaction) {
		return "transaction failed"
	}

	return "error occurred"
}
