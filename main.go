package main

import (
	"Library-Archive/api"
	"Library-Archive/db"
	"Library-Archive/handler"

	"github.com/labstack/echo/v4"
)

func main() {
	db.Connect()

	db.DB.Exec(`CREATE TABLE IF NOT EXISTS books (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        title TEXT,
        author_id INTEGER,
        published_year INTEGER
    )`)

	db.DB.Exec(`CREATE TABLE IF NOT EXISTS authors (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        name TEXT,
        bio TEXT
    )`)

	h := &handler.Handler{DB: db.DB}

	e := echo.New()
	api.RegisterHandlers(e, h)
	_ = e.Start(":8080")
}
