package models

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
)

type Models struct {
	PagesM PageModel
	UsersM UserModel
	TilsM  TilModel
	BolgM  BlogModel
}

func NewModel(db *sql.DB) Models {
	return Models{
		PagesM: PageModel{DB: db},
		UsersM: UserModel{DB: db},
		TilsM:  TilModel{DB: db},
		BolgM:  BlogModel{DB: db},
	}
}
