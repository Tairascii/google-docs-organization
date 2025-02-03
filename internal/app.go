package app

import "github.com/Tairascii/google-docs-organization/internal/app/usecase"

type UseCase struct {
	Org usecase.OrgUseCase
}

type DI struct {
	UseCase UseCase
}
