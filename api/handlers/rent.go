package handlers

import (
	"go-fiber-digilib/api/middlewares"
	"go-fiber-digilib/pkg/constant"
	"go-fiber-digilib/pkg/dtos"
	"go-fiber-digilib/pkg/rajaongkir"
	"go-fiber-digilib/pkg/utils"
	"go-fiber-digilib/rents"
	"go-fiber-digilib/settings"
	"go-fiber-digilib/users"
	"net/http"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v3"
)

type Rents struct {
	rents     rents.Service
	settings  settings.Service
	users     users.Service
	roService rajaongkir.Service
}

func (r Rents) GetAll(ctx fiber.Ctx) error {
	page, _ := strconv.Atoi(ctx.Query("page"))
	limit, _ := strconv.Atoi(ctx.Query("limit"))
	sort := ctx.Query("sort")
	search := ctx.Query("search")

	pagination := utils.Pagination{
		Page:    page,
		Limit:   limit,
		Sort:    sort,
		Search:  search,
		Keyword: "users.username",
	}

	booksData, err := r.rents.GetAll(ctx.Context(), pagination)

	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(dtos.Response[any]{
			Status:  "failed",
			Message: utils.GetErrorMessage(err),
		})
	}

	return ctx.Status(http.StatusOK).JSON(dtos.Response[utils.Pagination]{
		Status:  "success",
		Message: "all book rents",
		Data:    booksData,
	})
}

func (r Rents) GetByUser(ctx fiber.Ctx) error {
	userData := ctx.Locals("userData").(*middlewares.JWTCustomClaims)
	userID := userData.ID

	rentsData, err := r.rents.GetByUser(ctx.Context(), uint(userID))

	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(dtos.Response[any]{
			Status:  "failed",
			Message: utils.GetErrorMessage(err),
		})
	}

	return ctx.Status(http.StatusOK).JSON(dtos.Response[[]rents.UserRent]{
		Status:  "success",
		Message: "all rents",
		Data:    rentsData,
	})
}

func (r Rents) Create(ctx fiber.Ctx) error {
	userData := ctx.Locals("userData").(*middlewares.JWTCustomClaims)
	userID := userData.ID

	userRecord, err := r.users.GetProfile(ctx.Context(), uint(userID))

	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(dtos.Response[any]{
			Status:  "failed",
			Message: utils.GetErrorMessage(err),
		})
	}

	rentReq := ctx.Locals("validatedBody").(*rents.RentRequest)
	rentReq.UserID = uint(userID)

	setting, err := r.settings.GetByKey(ctx.Context(), constant.DISTRICT_ID)

	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(dtos.Response[any]{
			Status:  "failed",
			Message: utils.GetErrorMessage(err),
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
		return ctx.Status(http.StatusInternalServerError).JSON(dtos.Response[any]{
			Status:  "failed",
			Message: "failed to calculate fee",
		})
	}

	rentReq.ReturnTime = time.Now().AddDate(0, 0, int(rentReq.Duration))
	rentReq.Fee = fee

	rent, err := r.rents.Create(ctx.Context(), rentReq)

	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(dtos.Response[any]{
			Status:  "failed",
			Message: utils.GetErrorMessage(err),
		})
	}

	return ctx.Status(http.StatusCreated).JSON(dtos.Response[rents.Rent]{
		Status:  "success",
		Message: "book added to the rent",
		Data:    rent,
	})
}

func (r Rents) Update(ctx fiber.Ctx) error {
	param := ctx.Params("id")

	id, err := strconv.ParseUint(param, 10, 64)

	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(dtos.Response[any]{
			Status:  "failed",
			Message: "invalid id",
		})
	}

	rentReq := ctx.Locals("validatedBody").(*rents.RentUpdateRequest)

	rent, err := r.rents.Update(ctx.Context(), rentReq, uint(id))

	if err != nil {
		return ctx.Status(http.StatusNotFound).JSON(dtos.Response[any]{
			Status:  "failed",
			Message: utils.GetErrorMessage(err),
		})
	}

	return ctx.Status(http.StatusOK).JSON(dtos.Response[rents.Rent]{
		Status:  "success",
		Message: "rent updated",
		Data:    rent,
	})
}

func (r Rents) Delete(ctx fiber.Ctx) error {
	param := ctx.Params("id")

	id, err := strconv.ParseUint(param, 10, 64)

	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(dtos.Response[any]{
			Status:  "failed",
			Message: "invalid id",
		})
	}

	err = r.rents.Delete(ctx.Context(), uint(id))

	if err != nil {
		return ctx.Status(http.StatusNotFound).JSON(dtos.Response[any]{
			Status:  "failed",
			Message: utils.GetErrorMessage(err),
		})
	}

	return ctx.Status(http.StatusOK).JSON(dtos.Response[any]{
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
