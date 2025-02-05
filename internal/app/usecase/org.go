package usecase

import (
	"context"
	"errors"
	"github.com/Tairascii/google-docs-organization/internal/app/service/org"
	"github.com/google/uuid"
)

var (
	ErrInvalidOwnerId = errors.New("invalid owner id")
)

type OrgUseCase interface {
	CreateOrg(ctx context.Context, orgName string) (uuid.UUID, error)
	AddUser(ctx context.Context, orgId, userId uuid.UUID, role string) error
}

type UseCase struct {
	org org.OrgService
}

func NewOrgUseCase(orgSrv org.OrgService) *UseCase {
	return &UseCase{org: orgSrv}
}

func (u *UseCase) CreateOrg(ctx context.Context, orgName string) (uuid.UUID, error) {
	id, err := u.org.CreateOrg(ctx, orgName)
	if err != nil {
		if errors.Is(err, org.ErrInvalidOwnerId) {
			return uuid.Nil, ErrInvalidOwnerId
		}
		return uuid.Nil, err
	}
	return id, nil
}

func (u *UseCase) AddUser(ctx context.Context, orgId, userId uuid.UUID, role string) error {
	return u.org.AddUser(ctx, orgId, userId, role)
}
