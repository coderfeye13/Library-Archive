package db

import "gorm.io/gorm"

type Author struct {
	gorm.Model
	Name string
	Bio  string
}

type Book struct {
	gorm.Model
	Title         string
	PublishedYear int
	AuthorID      uint   // foreign key — Author.ID'ye işaret eder
	Author        Author // ilişki nesnesi — Preload ile dolar
}
