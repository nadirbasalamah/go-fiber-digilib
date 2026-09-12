package handlers

import (
	"go-digilib/pkg/dtos"
	"go-digilib/settings"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v5"
)

type Settings struct {
	settings settings.Service
}

func (s Settings) GetAll(ctx *echo.Context) error {

	settingsData, err := s.settings.GetAll(ctx.Request().Context())

	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, dtos.Response[any]{
			Status:  "failed",
			Message: "fetch settings failed",
		})
	}

	return ctx.JSON(http.StatusOK, dtos.Response[[]settings.Setting]{
		Status:  "success",
		Message: "all settings",
		Data:    settingsData,
	})
}

func (s Settings) GetByID(ctx *echo.Context) error {
	param := ctx.Param("id")

	id, err := strconv.ParseUint(param, 10, 64)

	if err != nil {
		return ctx.JSON(http.StatusBadRequest, dtos.Response[any]{
			Status:  "failed",
			Message: "invalid id",
		})
	}

	setting, err := s.settings.GetByID(ctx.Request().Context(), uint(id))

	if err != nil {
		return ctx.JSON(http.StatusNotFound, dtos.Response[any]{
			Status:  "failed",
			Message: "setting not found",
		})
	}

	return ctx.JSON(http.StatusOK, dtos.Response[settings.Setting]{
		Status:  "success",
		Message: "setting found",
		Data:    setting,
	})
}

func (s Settings) GetByKey(ctx *echo.Context) error {
	param := ctx.Param("key")

	setting, err := s.settings.GetByKey(ctx.Request().Context(), param)

	if err != nil {
		return ctx.JSON(http.StatusNotFound, dtos.Response[any]{
			Status:  "failed",
			Message: "setting not found",
		})
	}

	return ctx.JSON(http.StatusOK, dtos.Response[settings.Setting]{
		Status:  "success",
		Message: "setting found",
		Data:    setting,
	})
}

func (s Settings) Create(ctx *echo.Context) error {
	settingReq := ctx.Get("validatedBody").(*settings.SettingRequest)

	setting, err := s.settings.Create(ctx.Request().Context(), settingReq)

	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, dtos.Response[any]{
			Status:  "failed",
			Message: "create setting failed",
		})
	}

	return ctx.JSON(http.StatusCreated, dtos.Response[settings.Setting]{
		Status:  "success",
		Message: "setting created",
		Data:    setting,
	})
}

func (s Settings) Update(ctx *echo.Context) error {
	param := ctx.Param("id")

	id, err := strconv.ParseUint(param, 10, 64)

	if err != nil {
		return ctx.JSON(http.StatusBadRequest, dtos.Response[any]{
			Status:  "failed",
			Message: "invalid id",
		})
	}

	settingReq := ctx.Get("validatedBody").(*settings.SettingRequest)

	setting, err := s.settings.Update(ctx.Request().Context(), settingReq, uint(id))

	if err != nil {
		return ctx.JSON(http.StatusNotFound, dtos.Response[any]{
			Status:  "failed",
			Message: "update setting failed",
		})
	}

	return ctx.JSON(http.StatusOK, dtos.Response[settings.Setting]{
		Status:  "success",
		Message: "setting updated",
		Data:    setting,
	})
}

func (s Settings) Delete(ctx *echo.Context) error {
	param := ctx.Param("id")

	id, err := strconv.ParseUint(param, 10, 64)

	if err != nil {
		return ctx.JSON(http.StatusBadRequest, dtos.Response[any]{
			Status:  "failed",
			Message: "invalid id",
		})
	}

	err = s.settings.Delete(ctx.Request().Context(), uint(id))

	if err != nil {
		return ctx.JSON(http.StatusNotFound, dtos.Response[any]{
			Status:  "failed",
			Message: "setting not found",
		})
	}

	return ctx.JSON(http.StatusOK, dtos.Response[any]{
		Status:  "success",
		Message: "setting deleted",
	})
}

func NewSettings(settings settings.Service) Settings {
	settingsHandler := Settings{
		settings: settings,
	}

	return settingsHandler
}
