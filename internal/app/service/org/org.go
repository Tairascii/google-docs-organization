package org

import (
	"context"
	"errors"
	"github.com/Tairascii/google-docs-organization/internal/app/service/org/repo"
	"github.com/google/uuid"
)

var (
	ErrOnCreate       = errors.New("create error")
	ErrOnAdd          = errors.New("add error")
	ErrInvalidOwnerId = errors.New("invalid owner id")
)

const (
	adminRole = "admin"
)

type OrgService interface {
	CreateOrg(ctx context.Context, orgName string) (uuid.UUID, error)
}

type Service struct {
	repo repo.OrgRepo
}

func New(r repo.OrgRepo) *Service {
	return &Service{repo: r}
}

func (s *Service) CreateOrg(ctx context.Context, orgName string) (uuid.UUID, error) {
	ownerIdRaw, ok := ctx.Value("id").(string)
	if !ok {
		return uuid.Nil, ErrInvalidOwnerId
	}

	ownerId, err := uuid.Parse(ownerIdRaw)
	if err != nil {
		return uuid.Nil, ErrInvalidOwnerId
	}
	//TODO query as tx with adding user
	id, err := s.repo.CreateOrg(ctx, orgName)
	if err != nil {
		return uuid.Nil, errors.Join(ErrOnCreate, err)
	}

	err = s.repo.AddUserToOrg(ctx, id, ownerId, adminRole)
	if err != nil {
		return uuid.Nil, errors.Join(ErrOnAdd, err)
	}
	return id, nil
}
