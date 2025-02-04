package usecase

import (
	"context"
	"github.com/Tairascii/google-docs-organization/internal/app/service/org"
)

type OrgUseCase interface {
	CreateOrg(ctx context.Context, orgName string) (string, error)
}

type UseCase struct {
	org org.OrgService
}

func NewOrgUseCase(orgSrv org.OrgService) *UseCase {
	return &UseCase{org: orgSrv}
}

func (u *UseCase) CreateOrg(ctx context.Context, orgName string) (string, error) {
	return u.org.CreateOrg(ctx, orgName)
}
