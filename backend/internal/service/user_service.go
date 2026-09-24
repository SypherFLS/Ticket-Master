package service

import (
	"context"
	"tmaster/internal/apperrors"
	"tmaster/internal/auth"
	"tmaster/internal/constants"
	"tmaster/internal/dto"
)

func (s *Service) RegisterService(ctx context.Context, user dto.RegisterDTO) error {
	hashed, err := auth.HashPassword(user.Password)
	if err != nil {
		return err
	}

	userModel := dto.RegToModel(user, hashed)

	return s.repo.RegisterRepo(ctx, userModel)
}

func (s *Service) LoginService(ctx context.Context, user dto.LoginDTO) (string, error) {
	result, err := s.repo.GetLoginDataRepo(ctx, user.Email)

	if err != nil {
		return "", err
	}

	if err := auth.CheckPassword(user.Password, result.PasswordHash); err != nil {
		return "", err
	}

	token, err := s.jwtmanager.Generate(result.ID, result.UserRole)

	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *Service) SetRoleService(ctx context.Context, request_role constants.UserRole, email string, role_to_upd constants.UserRole) error {
	if !auth.HasPermission(request_role, auth.PermissionUserManage) {
		return apperrors.NoPermission
	}

	return s.repo.SetRoleRepo(ctx, email, role_to_upd)
}
