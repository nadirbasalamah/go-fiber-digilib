package handlers

import (
	"context"
	"go-fiber-digilib/api/middlewares"
	"go-fiber-digilib/carts"
	"go-fiber-digilib/pkg/dtos"
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v3"
)

type Carts struct {
	carts carts.Service
}

func (c Carts) GetByUser(ctx fiber.Ctx) error {
	userData := ctx.Locals("userData").(*middlewares.JWTCustomClaims)
	userID := userData.ID

	cartsData, err := c.carts.GetByUser(ctx.Context(), uint(userID))

	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(dtos.Response[any]{
			Status:  "failed",
			Message: "fetch carts failed",
		})
	}

	return ctx.Status(http.StatusOK).JSON(dtos.Response[[]carts.Cart]{
		Status:  "success",
		Message: "all carts",
		Data:    cartsData,
	})
}

func (c Carts) Create(ctx fiber.Ctx) error {
	userData := ctx.Locals("userData").(*middlewares.JWTCustomClaims)
	userID := userData.ID

	cartReq := ctx.Locals("validatedBody").(*carts.CartRequest)
	cartReq.UserID = uint(userID)

	cart, err := c.carts.Create(ctx.Context(), cartReq)

	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(dtos.Response[any]{
			Status:  "failed",
			Message: "add book to the cart failed",
		})
	}

	return ctx.Status(http.StatusCreated).JSON(dtos.Response[carts.Cart]{
		Status:  "success",
		Message: "book added to the cart",
		Data:    cart,
	})
}

func (c Carts) Update(ctx fiber.Ctx) error {
	param := ctx.Params("id")

	id, err := strconv.ParseUint(param, 10, 64)

	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(dtos.Response[any]{
			Status:  "failed",
			Message: "invalid id",
		})
	}

	userData := ctx.Locals("userData").(*middlewares.JWTCustomClaims)
	userID := userData.ID

	cartReq := ctx.Locals("validatedBody").(*carts.CartRequest)
	cartReq.UserID = uint(userID)

	cart, err := c.carts.Update(ctx.Context(), cartReq, uint(id))

	if err != nil {
		return ctx.Status(http.StatusNotFound).JSON(dtos.Response[any]{
			Status:  "failed",
			Message: "update cart failed",
		})
	}

	return ctx.Status(http.StatusOK).JSON(dtos.Response[carts.Cart]{
		Status:  "success",
		Message: "cart updated",
		Data:    cart,
	})
}

func (c Carts) Delete(ctx fiber.Ctx) error {
	param := ctx.Params("id")

	id, err := strconv.ParseUint(param, 10, 64)

	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(dtos.Response[any]{
			Status:  "failed",
			Message: "invalid id",
		})
	}

	userData := ctx.Locals("userData").(*middlewares.JWTCustomClaims)
	userID := userData.ID

	cartsCtx := context.WithValue(ctx.Context(), "userID", userID)

	err = c.carts.Delete(cartsCtx, uint(id))

	if err != nil {
		return ctx.Status(http.StatusNotFound).JSON(dtos.Response[any]{
			Status:  "failed",
			Message: "cart not found",
		})
	}

	return ctx.Status(http.StatusOK).JSON(dtos.Response[any]{
		Status:  "success",
		Message: "book removed from the cart",
	})
}

func NewCarts(carts carts.Service) Carts {
	cartsHandler := Carts{
		carts: carts,
	}

	return cartsHandler
}
