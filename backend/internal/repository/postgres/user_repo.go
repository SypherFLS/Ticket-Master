package postgres

import (
	"context"
	"tmaster/internal/apperrors"
	"tmaster/internal/constants"
	"tmaster/internal/repository/models"
)

func (r *Repo) RegisterRepo(ctx context.Context, user models.User) error {
	return r.db.WithContext(ctx).Model(&models.User{}).Create(user).Error
}

func (r *Repo) GetLoginDataRepo(ctx context.Context, email string) (models.LoginResult, error) {
	var data models.LoginResult
	res := r.db.WithContext(ctx).
		Model(models.User{}).
		Select("id", "password_hash", "role").
		Where("email = ?", email).
		First(&data).Error
	return data, res
}

func (r *Repo) SetRoleRepo(ctx context.Context,email string, role constants.UserRole) error { // rework jwt token
	res := r.db.WithContext(ctx).
		Model(&models.User{}).
		Where("email = ?", email).
		Update("role", role)

	if res.RowsAffected == 0 {
		return apperrors.NothingChanged
	}

	return res.Error
}

