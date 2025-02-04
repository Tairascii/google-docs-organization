package repo

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

var (
	ErrOnQuery = errors.New("query error")
)

type OrgRepo interface {
	CreateOrg(ctx context.Context, orgName string) (uuid.UUID, error)
	AddUserToOrg(ctx context.Context, orgID, userID uuid.UUID, role string) error
}

type Repo struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *Repo {
	return &Repo{db: db}
}

func (r *Repo) CreateOrg(ctx context.Context, orgName string) (uuid.UUID, error) {
	query := `INSERT INTO org (title) VALUES ($1) RETURNING id;`

	res := r.db.QueryRowxContext(ctx, query, orgName)

	var id uuid.UUID
	err := res.Scan(&id)
	if err != nil {
		return uuid.Nil, errors.Join(ErrOnQuery, err)
	}

	return id, nil
}

func (r *Repo) AddUserToOrg(ctx context.Context, orgID uuid.UUID, userID uuid.UUID, role string) error {
	query := `INSERT INTO org_members (org_id, user_id, role) VALUES ($1, $2, $3);`

	_, err := r.db.ExecContext(ctx, query, orgID, userID, role)
	if err != nil {
		return errors.Join(ErrOnQuery, err)
	}
	return nil
}
