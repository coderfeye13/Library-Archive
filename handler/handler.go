package handler

import (
	"Library-Archive/api"
	"Library-Archive/db"
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type Handler struct {
	DB *gorm.DB
}

func (h *Handler) ListBooks(ctx echo.Context) error {
	var books []db.Book
	// One call — GORM runs two SQL queries internally
	result := h.DB.Preload("Author").Find(&books)
	if result.Error != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": result.Error.Error(),
		})
	}
	// GORM runs and convert db.Book -> api.book before sending the response
	var response []api.Book
	for _, b := range books {
		id := int(b.ID)
		year := b.PublishedYear
		authorID := int(b.AuthorID)

		var apiAuthor *api.Author
		if b.Author.ID != 0 {
			authorObjID := int(b.Author.ID)
			apiAuthor = &api.Author{
				Id:   &authorObjID,
				Name: &b.Author.Name,
				Bio:  &b.Author.Bio,
			}
		}

		response = append(response, api.Book{
			Id:            &id,
			Title:         &b.Title,
			AuthorId:      &authorID,
			PublishedYear: &year,
			Author:        apiAuthor,
		})
	}

	return ctx.JSON(http.StatusOK, response)
}

func (h *Handler) CreateBook(ctx echo.Context) error {
	//API Layer clienttan gelen JSON oku
	var body api.NewBook

	if err := ctx.Bind(&body); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}
	//DB Layer GORM MOdel create and save
	book := db.Book{
		Title:    body.Title,
		AuthorID: uint(body.AuthorId),
	}
	if body.PublishedYear != nil {
		book.PublishedYear = *body.PublishedYear
	}

	result := h.DB.Create(&book)
	if result.Error != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": result.Error.Error(),
		})
	}
	//API Layer response api.Book
	id := int(book.ID)
	authorID := int(book.AuthorID)
	return ctx.JSON(http.StatusCreated, api.Book{
		Id:            &id,
		Title:         &book.Title,
		AuthorId:      &authorID,
		PublishedYear: &book.PublishedYear,
	})
}

func (h *Handler) GetBook(ctx echo.Context, id int) error {
	//DB Layer
	var book db.Book

	result := h.DB.Preload("Author").First(&book, id)
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
	//API Layer
	bookID := int(book.ID)
	authorID := int(book.AuthorID)

	var apiAuthor *api.Author
	if book.Author.ID != 0 {
		authorObjID := int(book.Author.ID)
		apiAuthor = &api.Author{
			Id:   &authorObjID,
			Name: &book.Author.Name,
			Bio:  &book.Author.Bio,
		}
	}

	return ctx.JSON(http.StatusOK, api.Book{
		Id:            &bookID,
		Title:         &book.Title,
		AuthorId:      &authorID,
		PublishedYear: &book.PublishedYear,
		Author:        apiAuthor,
	})
}

func (h *Handler) UpdateBook(ctx echo.Context, id int) error {
	var book db.Book
	var body api.NewBook

	if err := ctx.Bind(&body); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	result := h.DB.First(&book, id)
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

	book.Title = body.Title
	book.AuthorID = uint(body.AuthorId)
	if body.PublishedYear != nil {
		book.PublishedYear = *body.PublishedYear
	}

	result = h.DB.Save(&book)
	if result.Error != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": result.Error.Error(),
		})
	}

	bookID := int(book.ID)
	authorID := int(book.AuthorID)
	return ctx.JSON(http.StatusOK, api.Book{
		Id:            &bookID,
		Title:         &book.Title,
		AuthorId:      &authorID,
		PublishedYear: &book.PublishedYear,
	})
}

func (h *Handler) DeleteBook(ctx echo.Context, id int) error {
	result := h.DB.Delete(&db.Book{}, id)
	if result.Error != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": result.Error.Error(),
		})
	}

	return ctx.NoContent(http.StatusNoContent)
}

func (h *Handler) ListAuthors(ctx echo.Context) error {
	var authors []db.Author
	result := h.DB.Find(&authors)
	if result.Error != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": result.Error.Error(),
		})
	}

	var response []api.Author
	for _, a := range authors {
		id := int(a.ID)
		response = append(response, api.Author{
			Id:   &id,
			Name: &a.Name,
			Bio:  &a.Bio,
		})
	}

	return ctx.JSON(http.StatusOK, response)
}

func (h *Handler) CreateAuthor(ctx echo.Context) error {
	var body api.NewAuthor

	if err := ctx.Bind(&body); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	author := db.Author{
		Name: body.Name,
	}
	if body.Bio != nil {
		author.Bio = *body.Bio
	}

	result := h.DB.Create(&author)
	if result.Error != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": result.Error.Error(),
		})
	}

	id := int(author.ID)
	return ctx.JSON(http.StatusCreated, api.Author{
		Id:   &id,
		Name: &author.Name,
		Bio:  &author.Bio,
	})
}
