package data

import (
	"database/sql"
	"errors"
)

var (
	ErrorRecordNotFound = errors.New("record not found")
)

type Models struct {
	Movies interface {
		Insert(movie *Movie) error
		Get(id int64) (*Movie, error)
		Update(id int64) error
		Delete(movie *Movie) error
	}
}

func NewModels(db *sql.DB) Models {
	return Models{
		Movies: MovieModel{DB: db},
	}
}
