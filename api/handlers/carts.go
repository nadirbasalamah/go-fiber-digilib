package handlers

import (
	"context"
	"go-digilib/api/middlewares"
	"go-digilib/carts"
	"go-digilib/pkg/dtos"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v5"
)

type Carts struct {
	carts carts.Service
}

func (c Carts) GetByUser(ctx *echo.Context) error {
	userData := ctx.Get("userData").(*middlewares.JWTCustomClaims)
	userID := userData.ID

	cartsData, err := c.carts.GetByUser(ctx.Request().Context(), uint(userID))

	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, dtos.Response[any]{
			Status:  "failed",
			Message: "fetch carts failed",
		})
	}

	return ctx.JSON(http.StatusOK, dtos.Response[[]carts.Cart]{
		Status:  "success",
		Message: "all carts",
		Data:    cartsData,
	})
}

func (c Carts) Create(ctx *echo.Context) error {
	userData := ctx.Get("userData").(*middlewares.JWTCustomClaims)
	userID := userData.ID

	cartReq := ctx.Get("validatedBody").(*carts.CartRequest)
	cartReq.UserID = uint(userID)

	cart, err := c.carts.Create(ctx.Request().Context(), cartReq)

	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, dtos.Response[any]{
			Status:  "failed",
			Message: "add book to the cart failed",
		})
	}

	return ctx.JSON(http.StatusCreated, dtos.Response[carts.Cart]{
		Status:  "success",
		Message: "book added to the cart",
		Data:    cart,
	})
}

func (c Carts) Update(ctx *echo.Context) error {
	param := ctx.Param("id")

	id, err := strconv.ParseUint(param, 10, 64)

	if err != nil {
		return ctx.JSON(http.StatusBadRequest, dtos.Response[any]{
			Status:  "failed",
			Message: "invalid id",
		})
	}

	userData := ctx.Get("userData").(*middlewares.JWTCustomClaims)
	userID := userData.ID

	cartReq := ctx.Get("validatedBody").(*carts.CartRequest)
	cartReq.UserID = uint(userID)

	cart, err := c.carts.Update(ctx.Request().Context(), cartReq, uint(id))

	if err != nil {
		return ctx.JSON(http.StatusNotFound, dtos.Response[any]{
			Status:  "failed",
			Message: "update cart failed",
		})
	}

	return ctx.JSON(http.StatusOK, dtos.Response[carts.Cart]{
		Status:  "success",
		Message: "cart updated",
		Data:    cart,
	})
}

func (c Carts) Delete(ctx *echo.Context) error {
	param := ctx.Param("id")

	id, err := strconv.ParseUint(param, 10, 64)

	if err != nil {
		return ctx.JSON(http.StatusBadRequest, dtos.Response[any]{
			Status:  "failed",
			Message: "invalid id",
		})
	}

	userData := ctx.Get("userData").(*middlewares.JWTCustomClaims)
	userID := userData.ID

	userCtx := context.WithValue(ctx.Request().Context(), "userID", userID)
	req := ctx.Request().WithContext(userCtx)
	ctx.SetRequest(req)

	err = c.carts.Delete(ctx.Request().Context(), uint(id))

	if err != nil {
		return ctx.JSON(http.StatusNotFound, dtos.Response[any]{
			Status:  "failed",
			Message: "cart not found",
		})
	}

	return ctx.JSON(http.StatusOK, dtos.Response[any]{
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
