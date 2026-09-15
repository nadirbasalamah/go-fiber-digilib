package handlers

import (
	"go-fiber-digilib/pkg/dtos"
	"go-fiber-digilib/settings"
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v3"
)

type Settings struct {
	settings settings.Service
}

func (s Settings) GetAll(ctx fiber.Ctx) error {

	settingsData, err := s.settings.GetAll(ctx.Context())

	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(dtos.Response[any]{
			Status:  "failed",
			Message: "fetch settings failed",
		})
	}

	return ctx.Status(http.StatusOK).JSON(dtos.Response[[]settings.Setting]{
		Status:  "success",
		Message: "all settings",
		Data:    settingsData,
	})
}

func (s Settings) GetByID(ctx fiber.Ctx) error {
	param := ctx.Params("id")

	id, err := strconv.ParseUint(param, 10, 64)

	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(dtos.Response[any]{
			Status:  "failed",
			Message: "invalid id",
		})
	}

	setting, err := s.settings.GetByID(ctx.Context(), uint(id))

	if err != nil {
		return ctx.Status(http.StatusNotFound).JSON(dtos.Response[any]{
			Status:  "failed",
			Message: "setting not found",
		})
	}

	return ctx.Status(http.StatusOK).JSON(dtos.Response[settings.Setting]{
		Status:  "success",
		Message: "setting found",
		Data:    setting,
	})
}

func (s Settings) GetByKey(ctx fiber.Ctx) error {
	param := ctx.Params("key")

	setting, err := s.settings.GetByKey(ctx.Context(), param)

	if err != nil {
		return ctx.Status(http.StatusNotFound).JSON(dtos.Response[any]{
			Status:  "failed",
			Message: "setting not found",
		})
	}

	return ctx.Status(http.StatusOK).JSON(dtos.Response[settings.Setting]{
		Status:  "success",
		Message: "setting found",
		Data:    setting,
	})
}

func (s Settings) Create(ctx fiber.Ctx) error {
	settingReq := ctx.Locals("validatedBody").(*settings.SettingRequest)

	setting, err := s.settings.Create(ctx.Context(), settingReq)

	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(dtos.Response[any]{
			Status:  "failed",
			Message: "create setting failed",
		})
	}

	return ctx.Status(http.StatusCreated).JSON(dtos.Response[settings.Setting]{
		Status:  "success",
		Message: "setting created",
		Data:    setting,
	})
}

func (s Settings) Update(ctx fiber.Ctx) error {
	param := ctx.Params("id")

	id, err := strconv.ParseUint(param, 10, 64)

	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(dtos.Response[any]{
			Status:  "failed",
			Message: "invalid id",
		})
	}

	settingReq := ctx.Locals("validatedBody").(*settings.SettingRequest)

	setting, err := s.settings.Update(ctx.Context(), settingReq, uint(id))

	if err != nil {
		return ctx.Status(http.StatusNotFound).JSON(dtos.Response[any]{
			Status:  "failed",
			Message: "update setting failed",
		})
	}

	return ctx.Status(http.StatusOK).JSON(dtos.Response[settings.Setting]{
		Status:  "success",
		Message: "setting updated",
		Data:    setting,
	})
}

func (s Settings) Delete(ctx fiber.Ctx) error {
	param := ctx.Params("id")

	id, err := strconv.ParseUint(param, 10, 64)

	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(dtos.Response[any]{
			Status:  "failed",
			Message: "invalid id",
		})
	}

	err = s.settings.Delete(ctx.Context(), uint(id))

	if err != nil {
		return ctx.Status(http.StatusNotFound).JSON(dtos.Response[any]{
			Status:  "failed",
			Message: "setting not found",
		})
	}

	return ctx.Status(http.StatusOK).JSON(dtos.Response[any]{
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
