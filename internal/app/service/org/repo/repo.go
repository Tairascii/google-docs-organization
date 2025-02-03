package repo

import "github.com/jmoiron/sqlx"

type OrgRepo interface {
}

type Repo struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *Repo {
	return &Repo{db: db}
}
