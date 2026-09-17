package postgres

import (
	"context"
	"tmaster/internal/repository/models"
)

func (r *Repo) RegisterRepo(ctx context.Context, user models.User) error {
	return r.db.WithContext(ctx).Model(&models.User{}).Create(user).Error
}

func (r *Repo) GetLoginDataRepo(ctx context.Context, email string) (models.LoginResult, error) {
	var data models.LoginResult
	res := r.db.WithContext(ctx).Model(models.User{}).Select("id", "password_hash").Where("email = ?", email).First(&data)
	return data, res.Error
}

func (r *Repo) SetRoleRepo(ctx context.Context,email string, role models.UserRole) error {
	return r.db.WithContext(ctx).Model(&models.User{}).Where("email = ?", email).Update("role", role).Error
}
