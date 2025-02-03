package usecase

import "github.com/Tairascii/google-docs-organization/internal/app/service/org"

type OrgUseCase interface {
}

type UseCase struct {
	org org.OrgService
}

func NewOrgUseCase(orgSrv org.OrgService) *UseCase {
	return &UseCase{org: orgSrv}
}
