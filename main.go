package main

import (
	"Library-Archive/api"
	"Library-Archive/db"
	"Library-Archive/handler"
	"net/http"

	"github.com/labstack/echo/v4"
)

func main() {
	db.Connect()
	//GORM semayi yonetiyor otomatik olarak
	db.DB.AutoMigrate(&db.Author{}, &db.Book{}) // author once geliyor-> book tablosu foreign key pointing to authors

	h := &handler.Handler{DB: db.DB} //handler db yi alir

	e := echo.New()            //http server baslatir
	api.RegisterHandlers(e, h) //url-handler eslesme tanimlar
	//added for full specification over HTTP
	e.GET("/openapi.json", func(c echo.Context) error {
		swagger, err := api.GetSwagger()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": err.Error(),
			})
		}
		return c.JSON(http.StatusOK, swagger)
	})

	e.Static("/swagger", "static")

	e.Start(":8080")
}
