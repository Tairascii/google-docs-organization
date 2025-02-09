package usecase

import (
	"context"
	"errors"
	"github.com/Tairascii/google-docs-organization/internal/app/service/org"
	"github.com/Tairascii/google-docs-organization/internal/app/service/user"
	"github.com/google/uuid"
)

var (
	ErrInvalidOwnerId = errors.New("invalid owner id")
	ErrInvalidUserId  = errors.New("invalid user id")
)

type OrgUseCase interface {
	CreateOrg(ctx context.Context, orgName string) (uuid.UUID, error)
	AddUser(ctx context.Context, orgId uuid.UUID, email, role string) error
}

type UseCase struct {
	org org.OrgService
	usr user.UserService
}

func NewOrgUseCase(orgSrv org.OrgService, usrSrv user.UserService) *UseCase {
	return &UseCase{
		org: orgSrv,
		usr: usrSrv,
	}
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

func (u *UseCase) AddUser(ctx context.Context, orgId uuid.UUID, email, role string) error {
	userId, err := u.usr.IdByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, user.ErrNotFound) || errors.Is(err, user.ErrInvalidId) {
			return errors.Join(ErrInvalidUserId, err)
		}
		return err
	}

	return u.org.AddUser(ctx, orgId, userId, role)
}
