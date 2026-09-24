package service

import (
	"tmaster/internal/auth"
	"tmaster/internal/repository"
)

type Service struct {
	repo       repository.Repository
	jwtmanager *auth.JWTManager
}

func NewService(repo repository.Repository, jwtm *auth.JWTManager) *Service {
	return &Service{
		repo:       repo,
		jwtmanager: jwtm,
	}
}
