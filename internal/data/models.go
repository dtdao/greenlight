package data

import (
	"database/sql"
	"errors"
)

var (
	ErrorRecordNotFound = errors.New("record not found")
	ErrorEditConflict   = errors.New("edit conflict")
)

type Models struct {
	Users  UserModel
	Movies MovieModel
	Tokens TokenModel
	// this doesnt need to be an interface
	// Movies interface {
	// 	Insert(movie *Movie) error
	// 	Get(id int64) (*Movie, error)
	// 	Update(movie *Movie) error
	// 	Delete(id int64) error
	// 	GetAll(title string, genres []string, filters Filters) ([]*Movie, Metadata, error)
	// }
}

func NewModels(db *sql.DB) Models {
	return Models{
		Movies: MovieModel{DB: db},
		Users:  UserModel{DB: db},
		Tokens: TokenModel{DB: db},
	}
}
