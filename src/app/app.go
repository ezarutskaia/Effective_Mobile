package app

import (
	"test_effective_mobile/app/domain"
	"test_effective_mobile/app/repository"
	"test_effective_mobile/app/interfaces"
)

type App struct {
	Domain *domain.Domain
	Repository *repository.Repository
	Interfaces *interfaces.Interfaces
}