package repository

import (
	"context"

	"github.com/ngoctd314/c72-api-server/pkg/model"
)

func (r *Repository) CreateUser(ctx context.Context, user *model.User) error {
	tx := r.db.WithContext(ctx)

	if err := tx.Create(user).Error; err != nil {
		return err
	}

	return nil
}

func (r *Repository) GetUserByUsername(ctx context.Context, username string) (*model.User, error) {
	tx := r.db.WithContext(ctx)

	var user model.User
	if err := tx.Where(model.User{Username: username}).Take(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}
