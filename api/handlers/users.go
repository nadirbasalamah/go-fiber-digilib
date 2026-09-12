package handlers

import (
	"go-digilib/api/middlewares"
	"go-digilib/pkg/dtos"
	"go-digilib/pkg/fileupload"
	"go-digilib/pkg/utils"
	"go-digilib/users"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/labstack/echo/v5"
)

type Users struct {
	users users.Service
	cld   *cloudinary.Cloudinary
}

func (u Users) GetProfile(ctx *echo.Context) error {
	userID, err := middlewares.GetUserID(ctx.Request().Context())

	if err != nil {
		return ctx.JSON(http.StatusUnauthorized, dtos.Response[any]{
			Status:  "failed",
			Message: "invalid token",
		})
	}

	user, err := u.users.GetProfile(ctx.Request().Context(), uint(userID))

	if err != nil {
		return ctx.JSON(http.StatusNotFound, dtos.Response[any]{
			Status:  "failed",
			Message: "user not found",
		})
	}

	return ctx.JSON(http.StatusOK, dtos.Response[users.User]{
		Status:  "success",
		Message: "user data",
		Data:    user,
	})
}

func (u Users) EditProfile(ctx *echo.Context) error {
	userID, err := middlewares.GetUserID(ctx.Request().Context())

	if err != nil {
		return ctx.JSON(http.StatusUnauthorized, dtos.Response[any]{
			Status:  "failed",
			Message: "invalid token",
		})
	}

	editReq := new(users.EditProfileRequest)

	if err := ctx.Bind(editReq); err != nil {
		return ctx.JSON(http.StatusBadRequest, dtos.Response[any]{
			Status:  "failed",
			Message: "invalid request",
		})
	}

	if err := ctx.Validate(editReq); err != nil {
		return ctx.JSON(http.StatusUnprocessableEntity, dtos.Response[any]{
			Status:  "failed",
			Message: "validation failed",
			Data:    utils.GetValidationErrMessages(err.Error()),
		})
	}

	file, err := ctx.FormFile("file")
	if err != nil {
		return ctx.JSON(http.StatusUnprocessableEntity, dtos.Response[any]{
			Status:  "failed",
			Message: "file not found",
		})
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))

	isFileValid := utils.ValidateFile(ext)

	if !isFileValid {
		return ctx.JSON(http.StatusUnprocessableEntity, dtos.Response[any]{
			Status:  "failed",
			Message: "invalid file format",
		})
	}

	uploaderCfg := fileupload.CloudinaryUploader{
		Cld:     u.cld,
		File:    file,
		Options: uploader.UploadParams{},
	}

	fileLink, err := uploaderCfg.UploadFile(ctx.Request().Context())

	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, dtos.Response[any]{
			Status:  "failed",
			Message: "upload failed",
		})
	}

	editReq.ProfilePicture = fileLink

	user, err := u.users.Update(ctx.Request().Context(), editReq, uint(userID))

	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, dtos.Response[any]{
			Status:  "failed",
			Message: "update profile failed",
		})
	}

	return ctx.JSON(http.StatusOK, dtos.Response[users.User]{
		Status:  "success",
		Message: "profile updated",
		Data:    user,
	})
}

func NewUsers(users users.Service, cld *cloudinary.Cloudinary) Users {
	usersHandler := Users{
		users: users,
		cld:   cld,
	}

	return usersHandler
}
