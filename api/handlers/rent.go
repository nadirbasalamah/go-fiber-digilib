package handlers

import (
	"go-digilib/api/middlewares"
	"go-digilib/pkg/constant"
	"go-digilib/pkg/dtos"
	"go-digilib/pkg/rajaongkir"
	"go-digilib/pkg/utils"
	"go-digilib/rents"
	"go-digilib/settings"
	"go-digilib/users"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v5"
)

type Rents struct {
	rents     rents.Service
	settings  settings.Service
	users     users.Service
	roService rajaongkir.Service
}

func (r Rents) GetAll(ctx *echo.Context) error {
	page, _ := strconv.Atoi(ctx.QueryParam("page"))
	limit, _ := strconv.Atoi(ctx.QueryParam("limit"))
	sort := ctx.QueryParam("sort")
	search := ctx.QueryParam("search")

	pagination := utils.Pagination{
		Page:    page,
		Limit:   limit,
		Sort:    sort,
		Search:  search,
		Keyword: "users.username",
	}

	booksData, err := r.rents.GetAll(ctx.Request().Context(), pagination)

	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, dtos.Response[any]{
			Status:  "failed",
			Message: "fetch book rents failed",
		})
	}

	return ctx.JSON(http.StatusOK, dtos.Response[utils.Pagination]{
		Status:  "success",
		Message: "all book rents",
		Data:    booksData,
	})
}

func (r Rents) GetByUser(ctx *echo.Context) error {
	userData := ctx.Get("userData").(*middlewares.JWTCustomClaims)
	userID := userData.ID

	rentsData, err := r.rents.GetByUser(ctx.Request().Context(), uint(userID))

	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, dtos.Response[any]{
			Status:  "failed",
			Message: "fetch rents failed",
		})
	}

	return ctx.JSON(http.StatusOK, dtos.Response[[]rents.UserRent]{
		Status:  "success",
		Message: "all rents",
		Data:    rentsData,
	})
}

func (r Rents) Create(ctx *echo.Context) error {
	userData := ctx.Get("userData").(*middlewares.JWTCustomClaims)
	userID := userData.ID

	userRecord, err := r.users.GetProfile(ctx.Request().Context(), uint(userID))

	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, dtos.Response[any]{
			Status:  "failed",
			Message: "failed to retrieve user data",
		})
	}

	rentReq := ctx.Get("validatedBody").(*rents.RentRequest)
	rentReq.UserID = uint(userID)

	setting, err := r.settings.GetByKey(ctx.Request().Context(), constant.DISTRICT_ID)

	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, dtos.Response[any]{
			Status:  "failed",
			Message: "failed to retrieve origin",
		})
	}

	origin := setting.Value
	destination := strconv.Itoa(int(userRecord.DistrictID))

	fee, err := r.roService.GetDeliveryFee(rajaongkir.GetFeeRequest{
		Origin:      origin,
		Destination: destination,
		Courier:     rentReq.Courier,
	})

	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, dtos.Response[any]{
			Status:  "failed",
			Message: "failed to calculate fee",
		})
	}

	rentReq.ReturnTime = time.Now().AddDate(0, 0, int(rentReq.Duration))
	rentReq.Fee = fee

	rent, err := r.rents.Create(ctx.Request().Context(), rentReq)

	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, dtos.Response[any]{
			Status:  "failed",
			Message: "add book to the rent failed",
		})
	}

	return ctx.JSON(http.StatusCreated, dtos.Response[rents.Rent]{
		Status:  "success",
		Message: "book added to the rent",
		Data:    rent,
	})
}

func (r Rents) Update(ctx *echo.Context) error {
	param := ctx.Param("id")

	id, err := strconv.ParseUint(param, 10, 64)

	if err != nil {
		return ctx.JSON(http.StatusBadRequest, dtos.Response[any]{
			Status:  "failed",
			Message: "invalid id",
		})
	}

	rentReq := ctx.Get("validatedBody").(*rents.RentUpdateRequest)

	rent, err := r.rents.Update(ctx.Request().Context(), rentReq, uint(id))

	if err != nil {
		return ctx.JSON(http.StatusNotFound, dtos.Response[any]{
			Status:  "failed",
			Message: "update rent failed",
		})
	}

	return ctx.JSON(http.StatusOK, dtos.Response[rents.Rent]{
		Status:  "success",
		Message: "rent updated",
		Data:    rent,
	})
}

func (r Rents) Delete(ctx *echo.Context) error {
	param := ctx.Param("id")

	id, err := strconv.ParseUint(param, 10, 64)

	if err != nil {
		return ctx.JSON(http.StatusBadRequest, dtos.Response[any]{
			Status:  "failed",
			Message: "invalid id",
		})
	}

	err = r.rents.Delete(ctx.Request().Context(), uint(id))

	if err != nil {
		return ctx.JSON(http.StatusNotFound, dtos.Response[any]{
			Status:  "failed",
			Message: "rent not found",
		})
	}

	return ctx.JSON(http.StatusOK, dtos.Response[any]{
		Status:  "success",
		Message: "book rent removed",
	})
}

func NewRents(
	rents rents.Service,
	settings settings.Service,
	users users.Service,
	roService rajaongkir.Service,
) Rents {
	rentsHandler := Rents{
		rents:     rents,
		settings:  settings,
		users:     users,
		roService: roService,
	}

	return rentsHandler
}
