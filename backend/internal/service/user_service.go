package service

import (
	"context"
	"tmaster/internal/apperrors"
	"tmaster/internal/auth"
	"tmaster/internal/constants"
	"tmaster/internal/constants/params"
	"tmaster/internal/dto"
)

func (s *Service) RegisterService(ctx context.Context, user dto.RegisterDTO) error {
	hashed, err := auth.HashPassword(user.Password)
	if err != nil {
		return err
	}

	userModel := dto.RegToModel(user, hashed)

	if err := s.repo.RegisterRepo(ctx, userModel); err != nil {
		return err
	}

	logger := params.GetLogger(ctx)
	logger.Info(
		"user registered",
		"user_name", user.Name,
		"user_email", user.Email,
	)

	return nil
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

	err := s.repo.SetRoleRepo(ctx, email, role_to_upd)
	if err != nil {
		return err
	}

	logger := params.GetLogger(ctx)

	logger.Info(
		"user role changed",
		"user_email", email,
		"new_role", role_to_upd,
	)

	return nil
}
