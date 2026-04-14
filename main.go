package main

import (
	"Library-Archive/api"
	"Library-Archive/db"
	"Library-Archive/handler"
	"net/http"

	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func main() {
	db.Connect()
	db.DB.AutoMigrate(&db.Author{}, &db.Book{})

	h := &handler.Handler{DB: db.DB}

	e := echo.New()
	api.RegisterHandlers(e, h)

	e.GET("/openapi.json", func(c echo.Context) error {
		swagger, err := api.GetSwagger()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": err.Error(),
			})
		}
		return c.JSON(http.StatusOK, swagger)
	})

	e.GET("/swagger/*", echoSwagger.EchoWrapHandler(
		echoSwagger.URL("http://localhost:8080/openapi.json"),
	))

	e.Start(":8080")
}
