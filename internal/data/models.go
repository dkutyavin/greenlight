package data

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
)

func NewMoviesModel(db *sql.DB) *MovieModel {
	return &MovieModel{DB: db}
}
