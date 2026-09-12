package handlers

import (
	"encoding/json"
	"go-digilib/books"
	"go-digilib/pkg/ai"
	"go-digilib/pkg/dtos"
	"go-digilib/pkg/fileupload"
	"go-digilib/pkg/utils"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/labstack/echo/v5"
)

type Books struct {
	books          books.Service
	recommendation ai.Service
	cld            *cloudinary.Cloudinary
}

func (b Books) GetAll(ctx *echo.Context) error {
	page, _ := strconv.Atoi(ctx.QueryParam("page"))
	limit, _ := strconv.Atoi(ctx.QueryParam("limit"))
	sort := ctx.QueryParam("sort")
	search := ctx.QueryParam("search")

	pagination := utils.Pagination{
		Page:    page,
		Limit:   limit,
		Sort:    sort,
		Search:  search,
		Keyword: "title",
	}

	booksData, err := b.books.GetAll(ctx.Request().Context(), pagination)

	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, dtos.Response[any]{
			Status:  "failed",
			Message: "fetch books failed",
		})
	}

	return ctx.JSON(http.StatusOK, dtos.Response[utils.Pagination]{
		Status:  "success",
		Message: "all books",
		Data:    booksData,
	})
}

func (b Books) GetByID(ctx *echo.Context) error {
	param := ctx.Param("id")

	id, err := strconv.ParseUint(param, 10, 64)

	if err != nil {
		return ctx.JSON(http.StatusBadRequest, dtos.Response[any]{
			Status:  "failed",
			Message: "invalid id",
		})
	}

	book, err := b.books.GetByID(ctx.Request().Context(), uint(id))

	if err != nil {
		return ctx.JSON(http.StatusNotFound, dtos.Response[any]{
			Status:  "failed",
			Message: "book not found",
		})
	}

	return ctx.JSON(http.StatusOK, dtos.Response[books.Book]{
		Status:  "success",
		Message: "book found",
		Data:    book,
	})
}

func (b Books) GetByCategory(ctx *echo.Context) error {
	page, _ := strconv.Atoi(ctx.QueryParam("page"))
	limit, _ := strconv.Atoi(ctx.QueryParam("limit"))
	sort := ctx.QueryParam("sort")
	search := ctx.QueryParam("search")

	pagination := utils.Pagination{
		Page:    page,
		Limit:   limit,
		Sort:    sort,
		Search:  search,
		Keyword: "title",
	}

	param := ctx.Param("id")

	categoryId, err := strconv.ParseUint(param, 10, 64)

	if err != nil {
		return ctx.JSON(http.StatusBadRequest, dtos.Response[any]{
			Status:  "failed",
			Message: "invalid id",
		})
	}

	booksData, err := b.books.GetByCategory(ctx.Request().Context(), pagination, uint(categoryId))

	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, dtos.Response[any]{
			Status:  "failed",
			Message: "fetch books failed",
		})
	}

	return ctx.JSON(http.StatusOK, dtos.Response[utils.Pagination]{
		Status:  "success",
		Message: "all books by category",
		Data:    booksData,
	})
}

func (b Books) Create(ctx *echo.Context) error {
	bookReq := ctx.Get("validatedBody").(*books.BookRequest)

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
		Cld:     b.cld,
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

	bookReq.ImageLink = fileLink

	book, err := b.books.Create(ctx.Request().Context(), bookReq)

	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, dtos.Response[any]{
			Status:  "failed",
			Message: "create book failed",
		})
	}

	return ctx.JSON(http.StatusCreated, dtos.Response[books.Book]{
		Status:  "success",
		Message: "book created",
		Data:    book,
	})
}

func (b Books) Update(ctx *echo.Context) error {
	param := ctx.Param("id")

	id, err := strconv.ParseUint(param, 10, 64)

	if err != nil {
		return ctx.JSON(http.StatusBadRequest, dtos.Response[any]{
			Status:  "failed",
			Message: "invalid id",
		})
	}

	bookReq := ctx.Get("validatedBody").(*books.BookRequest)

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
		Cld:     b.cld,
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

	bookReq.ImageLink = fileLink

	book, err := b.books.Update(ctx.Request().Context(), bookReq, uint(id))

	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, dtos.Response[any]{
			Status:  "failed",
			Message: "update book failed",
		})
	}

	return ctx.JSON(http.StatusOK, dtos.Response[books.Book]{
		Status:  "success",
		Message: "book updated",
		Data:    book,
	})
}

func (b Books) Delete(ctx *echo.Context) error {
	param := ctx.Param("id")

	id, err := strconv.ParseUint(param, 10, 64)

	if err != nil {
		return ctx.JSON(http.StatusBadRequest, dtos.Response[any]{
			Status:  "failed",
			Message: "invalid id",
		})
	}

	err = b.books.Delete(ctx.Request().Context(), uint(id))

	if err != nil {
		return ctx.JSON(http.StatusNotFound, dtos.Response[any]{
			Status:  "failed",
			Message: "book not found",
		})
	}

	return ctx.JSON(http.StatusOK, dtos.Response[any]{
		Status:  "success",
		Message: "book deleted",
	})
}

func (b Books) GetBookRecommendation(ctx *echo.Context) error {
	recommReq := ctx.Get("validatedBody").(*ai.BookRecommendationRequest)

	req := ai.BookRecommendationRequest{
		Quantity: recommReq.Quantity,
		Topic:    recommReq.Topic,
	}

	res, err := b.recommendation.GetBookRecommendation(req)

	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, dtos.Response[any]{
			Status:  "failed",
			Message: "get book recommendation failed. Please retry soon",
		})
	}

	resBody := res.Choices[0].Message.Content

	var recommendations []ai.BookRecommendationResponse

	err = json.Unmarshal([]byte(resBody), &recommendations)

	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, dtos.Response[any]{
			Status:  "failed",
			Message: "get book recommendation failed. Invalid response",
		})
	}

	return ctx.JSON(http.StatusOK, dtos.Response[[]ai.BookRecommendationResponse]{
		Status:  "success",
		Message: "book recommendations",
		Data:    recommendations,
	})
}

func NewBooks(books books.Service, cld *cloudinary.Cloudinary, recommendation ai.Service) Books {
	booksHandler := Books{
		books:          books,
		cld:            cld,
		recommendation: recommendation,
	}

	return booksHandler
}
