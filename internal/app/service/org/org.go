package org

import "github.com/Tairascii/google-docs-organization/internal/app/service/org/repo"

type OrgService interface {
}

type Service struct {
	repo repo.OrgRepo
}

func New(r repo.OrgRepo) *Service {
	return &Service{repo: r}
}
