package utils

import "go-digilib/pkg/constant"

func ValidateFile(ext string) bool {
	return constant.ALLOWED_EXTENSIONS[ext]
}
