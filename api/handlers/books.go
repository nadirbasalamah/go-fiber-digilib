package handlers

import (
	"encoding/json"
	"go-fiber-digilib/books"
	"go-fiber-digilib/pkg/ai"
	"go-fiber-digilib/pkg/dtos"
	"go-fiber-digilib/pkg/fileupload"
	"go-fiber-digilib/pkg/utils"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v3"
)

type Books struct {
	books          books.Service
	recommendation ai.Service
	uploader       fileupload.Uploader
}

func (b Books) GetAll(ctx fiber.Ctx) error {
	page, _ := strconv.Atoi(ctx.Query("page"))
	limit, _ := strconv.Atoi(ctx.Query("limit"))
	sort := ctx.Query("sort")
	search := ctx.Query("search")

	pagination := utils.Pagination{
		Page:    page,
		Limit:   limit,
		Sort:    sort,
		Search:  search,
		Keyword: "title",
	}

	booksData, err := b.books.GetAll(ctx.Context(), pagination)

	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(dtos.Response[any]{
			Status:  "failed",
			Message: utils.GetErrorMessage(err),
		})
	}

	return ctx.Status(http.StatusOK).JSON(dtos.Response[utils.Pagination]{
		Status:  "success",
		Message: "all books",
		Data:    booksData,
	})
}

func (b Books) GetByID(ctx fiber.Ctx) error {
	param := ctx.Params("id")

	id, err := strconv.ParseUint(param, 10, 64)

	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(dtos.Response[any]{
			Status:  "failed",
			Message: "invalid id",
		})
	}

	book, err := b.books.GetByID(ctx.Context(), uint(id))

	if err != nil {
		return ctx.Status(http.StatusNotFound).JSON(dtos.Response[any]{
			Status:  "failed",
			Message: utils.GetErrorMessage(err),
		})
	}

	return ctx.Status(http.StatusOK).JSON(dtos.Response[books.Book]{
		Status:  "success",
		Message: "book found",
		Data:    book,
	})
}

func (b Books) GetByCategory(ctx fiber.Ctx) error {
	page, _ := strconv.Atoi(ctx.Query("page"))
	limit, _ := strconv.Atoi(ctx.Query("limit"))
	sort := ctx.Query("sort")
	search := ctx.Query("search")

	pagination := utils.Pagination{
		Page:    page,
		Limit:   limit,
		Sort:    sort,
		Search:  search,
		Keyword: "title",
	}

	param := ctx.Params("id")

	categoryId, err := strconv.ParseUint(param, 10, 64)

	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(dtos.Response[any]{
			Status:  "failed",
			Message: "invalid id",
		})
	}

	booksData, err := b.books.GetByCategory(ctx.Context(), pagination, uint(categoryId))

	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(dtos.Response[any]{
			Status:  "failed",
			Message: utils.GetErrorMessage(err),
		})
	}

	return ctx.Status(http.StatusOK).JSON(dtos.Response[utils.Pagination]{
		Status:  "success",
		Message: "all books by category",
		Data:    booksData,
	})
}

func (b Books) Create(ctx fiber.Ctx) error {
	bookReq := ctx.Locals("validatedBody").(*books.BookRequest)

	file, err := ctx.FormFile("file")
	if err != nil {
		return ctx.Status(http.StatusUnprocessableEntity).JSON(dtos.Response[any]{
			Status:  "failed",
			Message: "file not found",
		})
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))

	isFileValid := utils.ValidateFile(ext)

	if !isFileValid {
		return ctx.Status(http.StatusUnprocessableEntity).JSON(dtos.Response[any]{
			Status:  "failed",
			Message: "invalid file format",
		})
	}

	fileLink, err := b.uploader.UploadFile(ctx.Context(), file)

	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(dtos.Response[any]{
			Status:  "failed",
			Message: "upload failed",
		})
	}

	bookReq.ImageLink = fileLink

	book, err := b.books.Create(ctx.Context(), bookReq)

	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(dtos.Response[any]{
			Status:  "failed",
			Message: utils.GetErrorMessage(err),
		})
	}

	return ctx.Status(http.StatusCreated).JSON(dtos.Response[books.Book]{
		Status:  "success",
		Message: "book created",
		Data:    book,
	})
}

func (b Books) Update(ctx fiber.Ctx) error {
	param := ctx.Params("id")

	id, err := strconv.ParseUint(param, 10, 64)

	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(dtos.Response[any]{
			Status:  "failed",
			Message: "invalid id",
		})
	}

	bookReq := ctx.Locals("validatedBody").(*books.BookRequest)

	file, err := ctx.FormFile("file")
	if err != nil {
		return ctx.Status(http.StatusUnprocessableEntity).JSON(dtos.Response[any]{
			Status:  "failed",
			Message: "file not found",
		})
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))

	isFileValid := utils.ValidateFile(ext)

	if !isFileValid {
		return ctx.Status(http.StatusUnprocessableEntity).JSON(dtos.Response[any]{
			Status:  "failed",
			Message: "invalid file format",
		})
	}

	fileLink, err := b.uploader.UploadFile(ctx.Context(), file)

	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(dtos.Response[any]{
			Status:  "failed",
			Message: "upload failed",
		})
	}

	bookReq.ImageLink = fileLink

	book, err := b.books.Update(ctx.Context(), bookReq, uint(id))

	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(dtos.Response[any]{
			Status:  "failed",
			Message: utils.GetErrorMessage(err),
		})
	}

	return ctx.Status(http.StatusOK).JSON(dtos.Response[books.Book]{
		Status:  "success",
		Message: "book updated",
		Data:    book,
	})
}

func (b Books) Delete(ctx fiber.Ctx) error {
	param := ctx.Params("id")

	id, err := strconv.ParseUint(param, 10, 64)

	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(dtos.Response[any]{
			Status:  "failed",
			Message: "invalid id",
		})
	}

	err = b.books.Delete(ctx.Context(), uint(id))

	if err != nil {
		return ctx.Status(http.StatusNotFound).JSON(dtos.Response[any]{
			Status:  "failed",
			Message: utils.GetErrorMessage(err),
		})
	}

	return ctx.Status(http.StatusOK).JSON(dtos.Response[any]{
		Status:  "success",
		Message: "book deleted",
	})
}

func (b Books) GetBookRecommendation(ctx fiber.Ctx) error {
	recommReq := ctx.Locals("validatedBody").(*ai.BookRecommendationRequest)

	req := ai.BookRecommendationRequest{
		Quantity: recommReq.Quantity,
		Topic:    recommReq.Topic,
	}

	res, err := b.recommendation.GetBookRecommendation(req)

	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(dtos.Response[any]{
			Status:  "failed",
			Message: "get book recommendation failed. Please retry soon",
		})
	}

	resBody := res.Choices[0].Message.Content

	var recommendations []ai.BookRecommendationResponse

	err = json.Unmarshal([]byte(resBody), &recommendations)

	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(dtos.Response[any]{
			Status:  "failed",
			Message: "get book recommendation failed. Invalid response",
		})
	}

	return ctx.Status(http.StatusOK).JSON(dtos.Response[[]ai.BookRecommendationResponse]{
		Status:  "success",
		Message: "book recommendations",
		Data:    recommendations,
	})
}

func NewBooks(books books.Service, uploader fileupload.Uploader, recommendation ai.Service) Books {
	booksHandler := Books{
		books:          books,
		uploader:       uploader,
		recommendation: recommendation,
	}

	return booksHandler
}
