package handlers

import (
	"go-digilib/api/middlewares"
	"go-digilib/auth"
	"go-digilib/db/models"
	"go-digilib/pkg/dtos"
	"go-digilib/pkg/utils"
	"net/http"

	"github.com/labstack/echo/v5"
)

type Auth struct {
	auth      auth.Service
	jwtConfig middlewares.JWTConfig
}

func (a Auth) Register(ctx *echo.Context) error {
	registerReq := new(auth.RegisterRequest)

	if err := ctx.Bind(registerReq); err != nil {
		return ctx.JSON(http.StatusBadRequest, dtos.Response[any]{
			Status:  "failed",
			Message: "invalid request",
		})
	}

	if err := ctx.Validate(registerReq); err != nil {
		return ctx.JSON(http.StatusUnprocessableEntity, dtos.Response[any]{
			Status:  "failed",
			Message: "validation failed",
			Data:    utils.GetValidationErrMessages(err.Error()),
		})
	}

	user, err := a.auth.Register(ctx.Request().Context(), registerReq)

	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, dtos.Response[any]{
			Status:  "failed",
			Message: "user regsitration failed",
		})
	}

	return ctx.JSON(http.StatusCreated, dtos.Response[auth.User]{
		Status:  "success",
		Message: "user registered",
		Data:    user,
	})
}

func (a Auth) Login(ctx *echo.Context) error {
	loginReq := new(auth.LoginRequest)

	if err := ctx.Bind(loginReq); err != nil {
		return ctx.JSON(http.StatusBadRequest, dtos.Response[any]{
			Status:  "failed",
			Message: "invalid request",
		})
	}

	if err := ctx.Validate(loginReq); err != nil {
		return ctx.JSON(http.StatusUnprocessableEntity, dtos.Response[any]{
			Status:  "failed",
			Message: "validation failed",
			Data:    utils.GetValidationErrMessages(err.Error()),
		})
	}

	user, err := a.auth.Login(ctx.Request().Context(), loginReq)

	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, dtos.Response[any]{
			Status:  "failed",
			Message: "user login failed",
		})
	}

	token, err := a.jwtConfig.GenerateToken(int(user.ID), models.Role(user.Role))

	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, dtos.Response[any]{
			Status:  "failed",
			Message: "token generation failed",
		})
	}

	return ctx.JSON(http.StatusOK, dtos.Response[string]{
		Status:  "success",
		Message: "login success",
		Data:    token,
	})
}

func NewAuth(auth auth.Service, jwtConfig middlewares.JWTConfig) Auth {
	authHandler := Auth{
		auth:      auth,
		jwtConfig: jwtConfig,
	}

	return authHandler
}
