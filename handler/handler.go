package handler

import (
	"Library-Archive/api"
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type Handler struct {
	DB *gorm.DB
}

func (h *Handler) ListBooks(ctx echo.Context) error {
	var books []api.Book

	result := h.DB.Table("books").Find(&books)
	if result.Error != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": result.Error.Error(),
		})
	}

	return ctx.JSON(http.StatusOK, books)
}

func (h *Handler) CreateBook(ctx echo.Context) error {
	var body api.NewBook

	if err := ctx.Bind(&body); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	book := api.Book{
		Title:         &body.Title,
		AuthorId:      &body.AuthorId,
		PublishedYear: body.PublishedYear,
	}

	result := h.DB.Table("books").Create(&book)
	if result.Error != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": result.Error.Error(),
		})
	}

	return ctx.JSON(http.StatusCreated, book)
}

func (h *Handler) GetBook(ctx echo.Context, id int) error {
	var book api.Book

	result := h.DB.Table("books").First(&book, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return ctx.JSON(http.StatusNotFound, map[string]string{
				"error": "book not found",
			})
		}
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": result.Error.Error(),
		})
	}

	return ctx.JSON(http.StatusOK, book)
}

func (h *Handler) UpdateBook(ctx echo.Context, id int) error {
	var book api.Book
	var body api.NewBook

	if err := ctx.Bind(&body); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	result := h.DB.Table("books").First(&book, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return ctx.JSON(http.StatusNotFound, map[string]string{
				"error": "book not found",
			})
		}
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": result.Error.Error(),
		})
	}

	book.Title = &body.Title
	book.AuthorId = &body.AuthorId
	book.PublishedYear = body.PublishedYear

	result = h.DB.Table("books").Save(&book)
	if result.Error != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": result.Error.Error(),
		})
	}

	return ctx.JSON(http.StatusOK, book)
}

func (h *Handler) DeleteBook(ctx echo.Context, id int) error {
	result := h.DB.Table("books").Delete(&api.Book{}, id)
	if result.Error != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": result.Error.Error(),
		})
	}

	return ctx.NoContent(http.StatusNoContent)
}

func (h *Handler) ListAuthors(ctx echo.Context) error {
	var authors []api.Author

	result := h.DB.Table("authors").Find(&authors)
	if result.Error != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": result.Error.Error(),
		})
	}

	return ctx.JSON(http.StatusOK, authors)
}

func (h *Handler) CreateAuthor(ctx echo.Context) error {
	var body api.NewAuthor

	if err := ctx.Bind(&body); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	author := api.Author{
		Name: &body.Name,
		Bio:  body.Bio,
	}

	result := h.DB.Table("authors").Create(&author)
	if result.Error != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": result.Error.Error(),
		})
	}

	return ctx.JSON(http.StatusCreated, author)
}
