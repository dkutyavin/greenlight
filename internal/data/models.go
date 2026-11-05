package data

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")
)

func NewMoviesModel(db *sql.DB) *MovieModel {
	return &MovieModel{DB: db}
}

func NewUserModel(db *sql.DB) *UserModel {
	return &UserModel{DB: db}
}
