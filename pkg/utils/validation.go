package utils

import "go-fiber-digilib/pkg/constant"

func ValidateFile(ext string) bool {
	return constant.ALLOWED_EXTENSIONS[ext]
}
